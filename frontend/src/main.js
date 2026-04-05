/* Atualiza GO - Lógica Principal do Frontend */
/* Erasmo Cardoso - Software Engineer | Electronics Specialist */

import { GetDistroInfo, GetUpdateSteps, RunUpdate, RunSystemAction, GetSystemStats, IsRestrictedSandbox } from '../wailsjs/go/main/App';
import { EventsOn, BrowserOpenURL } from '../wailsjs/runtime/runtime';
import { translations } from './translations';
import './style.css';

const state = {
    distro: null,
    steps: [],
    selectedSteps: new Set(),
    running: false,
    language: localStorage.getItem('lang') || 'pt-br',
};

document.addEventListener('DOMContentLoaded', init);

async function init() {
    setupNavigation();
    setupActionButtons();
    setupLanguageToggle();
    applyTranslations();
    await loadDistroInfo();
    await checkSandboxMode();
    await loadUpdateSteps();
    listenProgress();
    
    // Inicia telemetria
    updateTelemetry();
    setInterval(updateTelemetry, 5000);
}

function t(key) {
    const lang = state.language;
    return translations[lang][key] || key;
}

function applyTranslations() {
    document.querySelectorAll('[data-i18n]').forEach(el => {
        const key = el.dataset.i18n;
        el.innerHTML = t(key);
    });

    // Atualiza placeholders dinâmicos
    if (state.distro) {
        document.getElementById('distro-name').textContent = state.distro.name || t('distro_unknown');
    }

    const btnUpdate = document.getElementById('btn-update');
    if (btnUpdate) {
        if (!state.running) {
            btnUpdate.innerHTML = `<span class="btn-icon">▶</span> ${t('btn_start_update')}`;
        } else {
            btnUpdate.innerHTML = `<span class="btn-icon">⏳</span> ${t('updating_label')}`;
        }
    }
}

async function checkSandboxMode() {
    try {
        const isSandbox = await IsRestrictedSandbox();
        if (isSandbox) {
            const banner = document.getElementById('sandbox-banner');
            if (banner) banner.classList.remove('hidden');
            
            // Trava o motor inteiro
            state.running = true;
            
            const btnUpdate = document.getElementById('btn-update');
            if (btnUpdate) {
                btnUpdate.disabled = true;
                btnUpdate.innerHTML = `<span class="btn-icon">🚫</span> Restrito`;
            }
        }
    } catch(err) {
        console.error('Sandbox detect error:', err);
    }
}

async function updateTelemetry() {
    try {
        const stats = await GetSystemStats();
        
        // RAM
        const ramBar = document.getElementById('ram-bar');
        const ramTxt = document.getElementById('ram-usage');
        if (ramBar && ramTxt) {
            const pct = Math.round(stats.memPercent);
            ramBar.style.width = pct + '%';
            ramBar.className = 'stats-bar-fill' + (pct > 90 ? ' danger' : pct > 75 ? ' warning' : '');
            ramTxt.textContent = `${stats.memUsed} MB / ${stats.memTotal} MB`;
        }

        // Disco
        const diskBar = document.getElementById('disk-bar');
        const diskTxt = document.getElementById('disk-usage');
        if (diskBar && diskTxt) {
            const pct = Math.round(stats.diskPercent);
            diskBar.style.width = pct + '%';
            diskBar.className = 'stats-bar-fill' + (pct > 90 ? ' danger' : pct > 75 ? ' warning' : '');
            diskTxt.textContent = `${stats.diskUsed} / ${stats.diskTotal}`;
        }
    } catch (err) {
        console.error('Telemetry error:', err);
    }
}

function setupLanguageToggle() {
    const btn = document.getElementById('btn-lang');
    if (!btn) return;
    btn.addEventListener('click', () => {
        state.language = state.language === 'pt-br' ? 'en' : 'pt-br';
        localStorage.setItem('lang', state.language);
        applyTranslations();
        renderUpdateOptions(state.steps);
    });
}

function setupNavigation() {
    document.querySelectorAll('.nav-item').forEach(btn => {
        btn.addEventListener('click', () => {
            const page = btn.dataset.page;
            document.querySelectorAll('.nav-item').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
            document.getElementById('page-' + page).classList.add('active');
        });
    });
}

async function loadDistroInfo() {
    try {
        const info = await GetDistroInfo();
        state.distro = info;

        document.getElementById('distro-name').textContent = info.name || t('distro_unknown');

        const icons = {
            debian: '🌀', ubuntu: '🟠', arch: '🔵',
            manjaro: '🟢', fedora: '🔵', mint: '🟢',
            opensuse: '🦎', alpine: '🏔️', void: '🔘', solus: '⛵'
        };
        document.getElementById('distro-icon').textContent = icons[info.id] || icons[info.family] || '🐧';

        document.getElementById('distro-family').textContent = info.family || '—';
        document.getElementById('distro-version').textContent = info.version ? 'v' + info.version : '—';
        document.getElementById('distro-pm').textContent = info.packageManager || '—';

        updateFeatureStatus('flatpak', info.hasFlatpak);
        updateFeatureStatus('snap', info.hasSnap);

    } catch (err) {
        document.getElementById('distro-name').textContent = 'Error';
    }
}

function updateFeatureStatus(feat, installed) {
    const statusEl = document.getElementById(`${feat}-status`);
    const actionsEl = document.getElementById(`${feat}-actions`);
    if (!statusEl || !actionsEl) return;

    if (installed) {
        statusEl.textContent = t('status_installed');
        statusEl.className = 'feature-status available';
        actionsEl.innerHTML = '';
    } else {
        statusEl.textContent = t('status_not_installed');
        statusEl.className = 'feature-status unavailable';
        actionsEl.innerHTML = `<button class="btn-install" data-feat="${feat}">${t('btn_install')}</button>`;
        
        actionsEl.querySelector('.btn-install').addEventListener('click', () => {
            // Vai para a página de atualização e seleciona o instalador
            document.querySelector('[data-page="update"]').click();
            const stepId = 'install_' + feat;
            if (state.steps.some(s => s.id === stepId)) {
                state.selectedSteps.add(stepId);
                renderUpdateOptions(state.steps);
                document.getElementById('btn-update').disabled = false;
            }
        });
    }
}

async function loadUpdateSteps() {
    try {
        const steps = await GetUpdateSteps();
        state.steps = steps;
        renderUpdateOptions(steps);
    } catch (err) {
        console.error('Error loading options:', err);
    }
}

function renderUpdateOptions(steps) {
    const container = document.getElementById('update-options');
    if (!container) return;
    container.innerHTML = '';

    if (steps.length === 0) {
        container.innerHTML = `<p style="color:var(--text-muted)">${t('update_no_options')}</p>`;
        return;
    }

    steps.forEach(step => {
        const el = document.createElement('div');
        el.className = 'update-option';
        if (state.selectedSteps.has(step.id)) el.classList.add('selected');
        el.dataset.stepId = step.id;
        
        const translatedLabel = t('step_' + step.id);
        const labelText = translatedLabel !== 'step_' + step.id ? translatedLabel : step.label;

        el.innerHTML = `
            <div class="update-checkbox"></div>
            <span class="update-option-label">${labelText}</span>
        `;
        el.addEventListener('click', () => toggleStep(step.id, el));
        container.appendChild(el);
    });
}

function toggleStep(id, el) {
    if (state.running) return;

    if (state.selectedSteps.has(id)) {
        state.selectedSteps.delete(id);
        el.classList.remove('selected');
    } else {
        state.selectedSteps.add(id);
        el.classList.add('selected');
    }

    document.getElementById('btn-update').disabled = state.selectedSteps.size === 0;
}

function setupActionButtons() {
    document.getElementById('btn-update').addEventListener('click', startUpdate);
    document.getElementById('btn-clear-log').addEventListener('click', () => {
        document.getElementById('log-output').innerHTML = '';
    });
    document.getElementById('btn-reboot').addEventListener('click', () => confirmAction('reboot'));
    document.getElementById('btn-shutdown').addEventListener('click', () => confirmAction('shutdown'));

    const btnDownload = document.getElementById('btn-sandbox-download');
    if (btnDownload) {
        btnDownload.addEventListener('click', () => {
            BrowserOpenURL('https://github.com/erascardsilva/atualiza_go/tree/main/build/bin');
        });
    }
}

async function startUpdate() {
    if (state.running || state.selectedSteps.size === 0) return;

    state.running = true;
    const btn = document.getElementById('btn-update');
    btn.disabled = true;
    btn.innerHTML = `<span class="btn-icon">⏳</span> ${t('updating_label')}`;

    const progressArea = document.getElementById('progress-area');
    progressArea.classList.remove('hidden');

    const stepsDiv = document.getElementById('progress-steps');
    stepsDiv.innerHTML = '';

    const ids = Array.from(state.selectedSteps);

    ids.forEach(id => {
        const step = state.steps.find(s => s.id === id);
        if (!step) return;

        const translatedLabel = t('step_' + id);
        const labelText = translatedLabel !== 'step_' + id ? translatedLabel : step.label;

        const el = document.createElement('div');
        el.className = 'progress-step';
        el.id = 'progress-' + id;
        el.innerHTML = `
            <div class="progress-step-header">
                <span class="progress-step-name">${labelText}</span>
                <span class="progress-step-pct">0%</span>
            </div>
            <div class="progress-bar-bg">
                <div class="progress-bar-fill" style="width:0%"></div>
            </div>
        `;
        stepsDiv.appendChild(el);
    });

    try {
        await RunUpdate(ids);
    } catch (err) {
        addLog('Error: ' + err, true);
    }
}

function listenProgress() {
    EventsOn('update:progress', (data) => {
        const el = document.getElementById('progress-' + data.stepId);
        if (el) {
            const pct = Math.min(100, Math.round(data.percent));
            el.querySelector('.progress-step-pct').textContent = pct + '%';
            el.querySelector('.progress-bar-fill').style.width = pct + '%';

            if (data.done) {
                el.classList.remove('running');
                el.classList.add(data.error ? 'error' : 'done');
            } else {
                el.classList.add('running');
            }
        }

        if (data.line) {
            addLog(data.line, !!data.error);
        }

        if (data.error) {
            addLog(data.error, true);
        }
    });

    EventsOn('update:complete', async () => {
        state.running = false;
        const btn = document.getElementById('btn-update');
        if (btn) {
            btn.innerHTML = `<span class="btn-icon">▶</span> ${t('btn_start_update')}`;
            btn.disabled = state.selectedSteps.size === 0;
        }

        const stepsDiv = document.getElementById('progress-steps');
        if (stepsDiv) {
            const banner = document.createElement('div');
            banner.className = 'complete-banner';
            banner.innerHTML = `
                <span class="complete-banner-icon">✅</span>
                <span class="complete-banner-text">${t('update_complete')}</span>
            `;
            stepsDiv.appendChild(banner);
        }
        
        // Recarrega info para ver se flatpak/snap foram instalados
        await loadDistroInfo();
        await loadUpdateSteps();
    });
}

function addLog(text, isError) {
    const log = document.getElementById('log-output');
    if (!log) return;
    const line = document.createElement('div');
    line.className = 'log-line' + (isError ? ' error' : '');
    line.textContent = text;
    log.appendChild(line);
    log.scrollTop = log.scrollHeight;
}

function confirmAction(action) {
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.innerHTML = `
        <div class="modal">
            <h3>${t(`modal_${action}_title`)}</h3>
            <p>${t(`modal_${action}_desc`)}</p>
            <div class="modal-actions">
                <button class="btn-cancel" id="modal-cancel">${t('modal_cancel')}</button>
                <button class="btn-danger" id="modal-confirm">${t('modal_confirm')}</button>
            </div>
        </div>
    `;

    document.body.appendChild(overlay);

    overlay.querySelector('#modal-cancel').addEventListener('click', () => overlay.remove());
    overlay.addEventListener('click', (e) => { if (e.target === overlay) overlay.remove(); });

    overlay.querySelector('#modal-confirm').addEventListener('click', async () => {
        overlay.remove();
        try {
            await RunSystemAction(action);
        } catch (err) {
            addLog('Error: ' + err, true);
        }
    });
}
