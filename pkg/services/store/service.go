package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grafana/grafana/pkg/infra/filestorage"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/registry"
	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/setting"
)

var grafanaStorageLogger = log.New("grafanaStorageLogger")

var ErrUploadFeatureDisabled = errors.New("upload feature is disabled")
var ErrUnsupportedStorage = errors.New("storage does not support this operation")
var ErrUploadInternalError = errors.New("upload internal error")
var ErrValidationFailed = errors.New("request validation failed")
var ErrFileAlreadyExists = errors.New("file exists")

const RootPublicStatic = "public-static"

const MAX_UPLOAD_SIZE = 1 * 1024 * 1024 // 3MB

type DeleteFolderCmd struct {
	Path  string `json:"path"`
	Force bool   `json:"force"`
}

type CreateFolderCmd struct {
	Path string `json:"path"`
}

type StorageService interface {
	registry.BackgroundService

	// List folder contents
	List(ctx context.Context, user *models.SignedInUser, path string) (*StorageListFrame, error)

	// Read raw file contents out of the store
	Read(ctx context.Context, user *models.SignedInUser, path string) (*filestorage.File, error)

	Upload(ctx context.Context, user *models.SignedInUser, req *UploadRequest) error

	Delete(ctx context.Context, user *models.SignedInUser, path string) error

	DeleteFolder(ctx context.Context, user *models.SignedInUser, cmd *DeleteFolderCmd) error

	CreateFolder(ctx context.Context, user *models.SignedInUser, cmd *CreateFolderCmd) error

	validateUploadRequest(ctx context.Context, user *models.SignedInUser, req *UploadRequest, storagePath string) validationResult

	// sanitizeUploadRequest sanitizes the upload request and converts it into a command accepted by the FileStorage API
	sanitizeUploadRequest(ctx context.Context, user *models.SignedInUser, req *UploadRequest, storagePath string) (*filestorage.UpsertFileCommand, error)
}

type storageServiceConfig struct {
	allowUnsanitizedSvgUpload bool
}

type standardStorageService struct {
	sql  *sqlstore.SQLStore
	tree *nestedTree
	cfg  storageServiceConfig
}

func ProvideService(sql *sqlstore.SQLStore, features featuremgmt.FeatureToggles, cfg *setting.Cfg) StorageService {
	settings, err := LoadStorageConfig(cfg)
	if err != nil {
		grafanaStorageLogger.Warn("error loading storage config", "error", err)
	}

	// always exists
	globalRoots := []storageRuntime{
		newDiskStorage(RootStorageConfig{
			Prefix:      RootPublicStatic,
			Name:        "Public static files",
			Description: "Access files from the static public files",
			Disk: &StorageLocalDiskConfig{
				Path: cfg.StaticRootPath,
				Roots: []string{
					"/testdata/",
					"/img/",
					"/gazetteer/",
					"/maps/",
				},
			},
		}).setReadOnly(true).setBuiltin(true),
	}

	// Development dashboards
	if settings.AddDevEnv && setting.Env != setting.Prod {
		devenv := filepath.Join(cfg.StaticRootPath, "..", "devenv")
		if _, err := os.Stat(devenv); !os.IsNotExist(err) {
			// path/to/whatever exists
			s := newDiskStorage(RootStorageConfig{
				Prefix:      "devenv",
				Name:        "Development Environment",
				Description: "Explore files within the developer environment directly",
				Disk: &StorageLocalDiskConfig{
					Path: devenv,
					Roots: []string{
						"/dev-dashboards/",
					},
				}}).setReadOnly(false)
			globalRoots = append(globalRoots, s)
		}
	}

	for _, root := range settings.Roots {
		if root.Prefix == "" {
			grafanaStorageLogger.Warn("Invalid root configuration", "cfg", root)
			continue
		}
		s, err := newStorage(root, filepath.Join(cfg.DataPath, "storage", "cache", root.Prefix))
		if err != nil {
			grafanaStorageLogger.Warn("error loading storage config", "error", err)
		}
		if s != nil {
			globalRoots = append(globalRoots, s)
		}
	}

	initializeOrgStorages := func(orgId int64) []storageRuntime {
		storages := make([]storageRuntime, 0)
		if features.IsEnabled(featuremgmt.FlagStorageLocalUpload) {
			storages = append(storages,
				newSQLStorage("resources",
					"Resources",
					"Upload custom resource files",
					&StorageSQLConfig{}, sql, orgId).
					setBuiltin(true))
		}
		return storages
	}

	return newStandardStorageService(sql, globalRoots, initializeOrgStorages)
}

func newStandardStorageService(sql *sqlstore.SQLStore, globalRoots []storageRuntime, initializeOrgStorages func(orgId int64) []storageRuntime) *standardStorageService {
	rootsByOrgId := make(map[int64][]storageRuntime)
	rootsByOrgId[ac.GlobalOrgID] = globalRoots

	res := &nestedTree{
		initializeOrgStorages: initializeOrgStorages,
		rootsByOrgId:          rootsByOrgId,
	}
	res.init()
	return &standardStorageService{
		sql:  sql,
		tree: res,
		cfg: storageServiceConfig{
			allowUnsanitizedSvgUpload: false,
		},
	}
}

func getOrgId(user *models.SignedInUser) int64 {
	if user == nil {
		return ac.GlobalOrgID
	}

	return user.OrgId
}

func (s *standardStorageService) Run(ctx context.Context) error {
	grafanaStorageLogger.Info("storage starting")
	return nil
}

func (s *standardStorageService) List(ctx context.Context, user *models.SignedInUser, path string) (*StorageListFrame, error) {
	// apply access control here

	return s.tree.ListFolder(ctx, getOrgId(user), path)
}

func (s *standardStorageService) Read(ctx context.Context, user *models.SignedInUser, path string) (*filestorage.File, error) {
	// TODO: permission check!
	return s.tree.GetFile(ctx, getOrgId(user), path)
}

type UploadRequest struct {
	Contents           []byte
	MimeType           string // TODO: remove MimeType from the struct once we can infer it from file contents
	Path               string
	CacheControl       string
	ContentDisposition string
	Properties         map[string]string
	EntityType         EntityType

	OverwriteExistingFile bool
}

func (s *standardStorageService) Upload(ctx context.Context, user *models.SignedInUser, req *UploadRequest) error {
	upload, storagePath := s.tree.getRoot(getOrgId(user), req.Path)
	if upload == nil {
		return ErrUploadFeatureDisabled
	}

	if upload.Meta().ReadOnly {
		return ErrUnsupportedStorage
	}

	validationResult := s.validateUploadRequest(ctx, user, req, storagePath)
	if !validationResult.ok {
		grafanaStorageLogger.Warn("file upload validation failed", "filetype", req.MimeType, "path", req.Path, "reason", validationResult.reason)
		return ErrValidationFailed
	}

	upsertCommand, err := s.sanitizeUploadRequest(ctx, user, req, storagePath)
	if err != nil {
		grafanaStorageLogger.Error("failed while sanitizing the upload request", "filetype", req.MimeType, "path", req.Path, "error", err)
		return ErrUploadInternalError
	}

	grafanaStorageLogger.Info("uploading a file", "filetype", req.MimeType, "path", req.Path)

	if !req.OverwriteExistingFile {
		file, err := upload.Store().Get(ctx, storagePath)
		if err != nil {
			grafanaStorageLogger.Error("failed while checking file existence", "err", err, "path", req.Path)
			return ErrUploadInternalError
		}

		if file != nil {
			return ErrFileAlreadyExists
		}
	}

	if err := upload.Store().Upsert(ctx, upsertCommand); err != nil {
		grafanaStorageLogger.Error("failed while uploading the file", "err", err, "path", req.Path)
		return ErrUploadInternalError
	}

	return nil
}

func (s *standardStorageService) DeleteFolder(ctx context.Context, user *models.SignedInUser, cmd *DeleteFolderCmd) error {
	root, storagePath := s.tree.getRoot(getOrgId(user), cmd.Path)
	if root == nil {
		return fmt.Errorf("resources storage is not enabled")
	}

	if root.Meta().ReadOnly {
		return ErrUnsupportedStorage
	}

	if storagePath == "" {
		storagePath = filestorage.Delimiter
	}
	return root.Store().DeleteFolder(ctx, storagePath, &filestorage.DeleteFolderOptions{Force: true})
}

func (s *standardStorageService) CreateFolder(ctx context.Context, user *models.SignedInUser, cmd *CreateFolderCmd) error {
	root, storagePath := s.tree.getRoot(getOrgId(user), cmd.Path)
	if root == nil {
		return fmt.Errorf("resources storage is not enabled")
	}

	if root.Meta().ReadOnly {
		return ErrUnsupportedStorage
	}

	err := root.Store().CreateFolder(ctx, storagePath)
	if err != nil {
		return err
	}
	return nil
}

func (s *standardStorageService) Delete(ctx context.Context, user *models.SignedInUser, path string) error {
	root, storagePath := s.tree.getRoot(getOrgId(user), path)
	if root == nil {
		return fmt.Errorf("resources storage is not enabled")
	}

	if root.Meta().ReadOnly {
		return ErrUnsupportedStorage
	}

	err := root.Store().Delete(ctx, storagePath)
	if err != nil {
		return err
	}
	return nil
}
