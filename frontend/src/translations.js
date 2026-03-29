/* Atualiza GO - Dicionário de Idiomas (i18n) */
/* Erasmo Cardoso - Dev */

export const translations = {
    'pt-br': {
        // Sidebar
        'nav_home': 'Início',
        'nav_update': 'Atualizar',
        'nav_actions': 'Ações',
        'nav_about': 'Sobre',

        // Home Page
        'home_title': 'Sistema Detectado',
        'home_subtitle': 'Informações do seu sistema',
        'distro_detecting': 'Detectando...',
        'distro_unknown': 'Desconhecido',
        'feat_flatpak': 'Flatpak',
        'feat_snap': 'Snap',
        'status_installed': 'Instalado',
        'status_up_to_date': 'Sistema Atualizado',

        // Update Page
        'update_title': 'Atualizar Sistema',
        'update_subtitle': 'Selecione as atualizações',
        'update_no_options': 'Nenhuma opção de atualização disponível para este sistema.',
        'btn_start_update': 'Iniciar Atualização',
        'updating_label': 'Atualizando...',
        'log_title': 'Log de Saída',
        'btn_clear_log': 'Limpar',
        'update_complete': 'Atualização concluída',

        // Actions Page
        'actions_title': 'Ações do Sistema',
        'actions_subtitle': 'Após a atualização',
        'action_reboot_label': 'Reiniciar',
        'action_reboot_desc': 'Reiniciar o sistema',
        'action_shutdown_label': 'Desligar',
        'action_shutdown_desc': 'Desligar o sistema',

        // About Page
        'about_title': 'Sobre',
        'about_subtitle': 'Atualiza GO — Atualizador de Sistema Linux',
        'how_to_use': 'Como Usar',
        'step_1': 'Na tela <strong>Início</strong>, confira se o sistema foi detectado corretamente.',
        'step_2': 'Acesse <strong>Atualizar</strong> e selecione os pacotes que deseja atualizar (Sistema, Flatpak, Snap).',
        'step_3': 'Clique em <strong>Iniciar Atualização</strong>. O sistema pedirá sua senha via autenticação gráfica.',
        'step_4': 'Acompanhe o progresso pelas barras e pelo log de saída em tempo real.',
        'step_5': 'Após concluir, use a tela <strong>Ações</strong> para reiniciar ou desligar, se necessário.',
        'distros_supported': 'Distros Suportadas',
        'distros_list': 'Debian, Ubuntu, Arch, Fedora, openSUSE e derivados.',
        'tech_title': 'Tecnologias',
        'tech_list': 'Go + Wails v2 ∙ HTML/CSS/JS ∙ pkexec (Polkit)',

        // Modals
        'modal_reboot_title': 'Reiniciar Sistema',
        'modal_reboot_desc': 'O sistema será reiniciado. Salve seu trabalho antes de continuar.',
        'modal_shutdown_title': 'Desligar Sistema',
        'modal_shutdown_desc': 'O sistema será desligado. Salve seu trabalho antes de continuar.',
        'modal_cancel': 'Cancelar',
        'modal_confirm': 'Confirmar',

        // Update Step Labels (Match IDs from Go)
        'step_system_update': 'Atualizar Sistema',
        'step_flatpak_update': 'Atualizar Flatpak',
        'step_snap_update': 'Atualizar Snap',
        'step_system_cleanup': 'Limpeza de Sistema',
        'step_install_flatpak': 'Instalar Suporte Flatpak',
        'step_install_snap': 'Instalar Suporte Snap',

        // System Stats
        'stats_ram': 'Memória RAM',
        'stats_disk': 'Disco Rígido (/)',
        'status_not_installed': 'Não Instalado',
        'btn_install': 'Instalar',
    },
    'en': {
        // Sidebar
        'nav_home': 'Home',
        'nav_update': 'Update',
        'nav_actions': 'Actions',
        'nav_about': 'About',

        // Home Page
        'home_title': 'System Detected',
        'home_subtitle': 'Your system information',
        'distro_detecting': 'Detecting...',
        'distro_unknown': 'Unknown',
        'feat_flatpak': 'Flatpak',
        'feat_snap': 'Snap',
        'status_installed': 'Installed',
        'status_not_installed': 'Not Installed',
        'status_up_to_date': 'System Updated',
        'btn_install': 'Install',

        // Update Page
        'update_title': 'Update System',
        'update_subtitle': 'Select the updates',
        'update_no_options': 'No update options available for this system.',
        'btn_start_update': 'Start Update',
        'updating_label': 'Updating...',
        'log_title': 'Output Log',
        'btn_clear_log': 'Clear',
        'update_complete': 'Update completed',

        // Actions Page
        'actions_title': 'System Actions',
        'actions_subtitle': 'After update',
        'action_reboot_label': 'Restart',
        'action_reboot_desc': 'Restart system',
        'action_shutdown_label': 'Shutdown',
        'action_shutdown_desc': 'Shutdown system',

        // About Page
        'about_title': 'About',
        'about_subtitle': 'Atualiza GO — Linux System Updater',
        'how_to_use': 'How to Use',
        'step_1': 'On the <strong>Home</strong> screen, check if the system was correctly detected.',
        'step_2': 'Go to <strong>Update</strong> and select the packages you want to update (System, Flatpak, Snap).',
        'step_3': 'Click <strong>Start Update</strong>. The system will ask for your password via graphical authentication.',
        'step_4': 'Track the progress through the bars and the real-time output log.',
        'step_5': 'After finishing, use the <strong>Actions</strong> screen to restart or shutdown, if necessary.',
        'distros_supported': 'Supported Distros',
        'distros_list': 'Debian, Ubuntu, Arch, Fedora, openSUSE and derivatives.',
        'tech_title': 'Technologies',
        'tech_list': 'Go + Wails v2 ∙ HTML/CSS/JS ∙ pkexec (Polkit)',

        // Modals
        'modal_reboot_title': 'Restart System',
        'modal_reboot_desc': 'The system will restart. Save your work before continuing.',
        'modal_shutdown_title': 'Shutdown System',
        'modal_shutdown_desc': 'The system will shutdown. Save your work before continuing.',
        'modal_cancel': 'Cancel',
        'modal_confirm': 'Confirm',

        // Update Step Labels
        'step_system_update': 'Update System',
        'step_flatpak_update': 'Update Flatpak',
        'step_snap_update': 'Update Snap',
        'step_system_cleanup': 'System Cleanup',
        'step_install_flatpak': 'Install Flatpak Support',
        'step_install_snap': 'Install Snap Support',

        // System Stats
        'stats_ram': 'RAM Memory',
        'stats_disk': 'Hard Drive (/)',
    }
}
