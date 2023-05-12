import { css } from '@emotion/css';
import React from 'react';

import { GrafanaTheme2, NavModelItem } from '@grafana/data';
import { Components } from '@grafana/e2e-selectors';
import { Dropdown, IconButton, ToolbarButton, useStyles2 } from '@grafana/ui';
import { contextSrv } from 'app/core/core';
import { t } from 'app/core/internationalization';
import { HOME_NAV_ID } from 'app/core/reducers/navModel';
import { useSelector } from 'app/types';

import { Breadcrumbs } from '../../Breadcrumbs/Breadcrumbs';
import { buildBreadcrumbs } from '../../Breadcrumbs/utils';
import { NewsContainer } from '../News/NewsContainer';
import { TopNavBarMenu } from '../TopBar/TopNavBarMenu';
import { TopSearchBarCommandPaletteTrigger } from '../TopBar/TopSearchBarCommandPaletteTrigger';
import { TOP_BAR_LEVEL_HEIGHT } from '../types';

export interface Props {
  onToggleSearchBar(): void;
  onToggleMegaMenu(): void;
  onToggleKioskMode(): void;
  searchBarHidden?: boolean;
  sectionNav: NavModelItem;
  pageNav?: NavModelItem;
  actions: React.ReactNode;
  megaMenuPinned?: boolean;
}

export function NavToolbar({
  actions,
  megaMenuPinned,
  searchBarHidden,
  sectionNav,
  pageNav,
  onToggleMegaMenu,
  onToggleSearchBar,
  onToggleKioskMode,
}: Props) {
  const homeNav = useSelector((state) => state.navIndex)[HOME_NAV_ID];
  const styles = useStyles2(getStyles);
  const breadcrumbs = buildBreadcrumbs(sectionNav, pageNav, homeNav);
  const navIndex = useSelector((state) => state.navIndex);
  const profileNode = navIndex['profile'];

  return (
    <div data-testid={Components.NavToolbar.container} className={styles.pageToolbar}>
      {!megaMenuPinned && (
        <div className={styles.menuButton}>
          <IconButton
            name="bars"
            tooltip={t('navigation.toolbar.toggle-menu', 'Toggle menu')}
            tooltipPlacement="bottom"
            size="xl"
            onClick={onToggleMegaMenu}
          />
        </div>
      )}
      <Breadcrumbs breadcrumbs={breadcrumbs} className={styles.breadcrumbsWrapper} />
      <div className={styles.actions}>
        {actions}
        {searchBarHidden && (
          <ToolbarButton
            onClick={onToggleKioskMode}
            narrow
            variant="default"
            title={t('navigation.toolbar.enable-kiosk', 'Enable kiosk mode')}
            icon="monitor"
          />
        )}
        <TopSearchBarCommandPaletteTrigger />
        <NewsContainer />
        {profileNode && (
          <Dropdown overlay={() => <TopNavBarMenu node={profileNode} />} placement="bottom-end">
            <ToolbarButton
              variant="default"
              className={styles.profileButton}
              imgSrc={contextSrv.user.gravatarUrl}
              imgAlt="User avatar"
              aria-label="Profile"
            />
          </Dropdown>
        )}
      </div>
    </div>
  );
}

const getStyles = (theme: GrafanaTheme2) => {
  return {
    breadcrumbsWrapper: css({
      display: 'flex',
      overflow: 'hidden',
      [theme.breakpoints.down('sm')]: {
        minWidth: '50%',
      },
    }),
    pageToolbar: css({
      height: TOP_BAR_LEVEL_HEIGHT,
      display: 'flex',
      padding: theme.spacing(0, 1, 0, 4),
      alignItems: 'center',
      flexShrink: 0,
      borderBottom: `1px solid ${theme.colors.border.weak}`,
    }),
    menuButton: css({
      display: 'flex',
      alignItems: 'center',
      marginRight: theme.spacing(1),
    }),
    logo: css({
      display: 'flex',
    }),
    img: css({
      height: theme.spacing(2),
      width: theme.spacing(2),
      marginRight: theme.spacing(1),
    }),
    actions: css({
      label: 'NavToolbar-actions',
      display: 'flex',
      alignItems: 'center',
      flexWrap: 'nowrap',
      justifyContent: 'flex-end',
      paddingLeft: theme.spacing(1),
      flexGrow: 1,
      gap: theme.spacing(0.5),
      minWidth: 0,

      '.body-drawer-open &': {
        display: 'none',
      },
    }),
    profileButton: css({
      padding: theme.spacing(0, 0.25),
      img: {
        borderRadius: '50%',
        height: '24px',
        marginRight: 0,
        width: '24px',
      },
    }),
  };
};
