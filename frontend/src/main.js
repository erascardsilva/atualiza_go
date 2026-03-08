import { GetDistroInfo, GetUpdateSteps, RunUpdate, RunSystemAction } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';
import './style.css';

const state = {
    distro: null,
    steps: [],
    selectedSteps: new Set(),
    running: false,
};

document.addEventListener('DOMContentLoaded', init);

async function init() {
    setupNavigation();
    setupActionButtons();
    await loadDistroInfo();
    await loadUpdateSteps();
    listenProgress();
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

        document.getElementById('distro-name').textContent = info.name || 'Desconhecido';

        const icons = {
            debian: '🟠', ubuntu: '🟠', arch: '🔵',
            manjaro: '🟢', fedora: '🔵', mint: '🟢',
        };
        document.getElementById('distro-icon').textContent = icons[info.id] || '🐧';

        document.getElementById('distro-family').textContent = info.family || '—';
        document.getElementById('distro-version').textContent = info.version ? 'v' + info.version : '—';
        document.getElementById('distro-pm').textContent = info.packageManager || '—';

        const flatStatus = document.getElementById('flatpak-status');
        const snapStatus = document.getElementById('snap-status');

        if (info.hasFlatpak) {
            flatStatus.textContent = 'Instalado';
            flatStatus.className = 'feature-status available';
        } else {
            flatStatus.textContent = 'Indisponível';
            flatStatus.className = 'feature-status unavailable';
        }

        if (info.hasSnap) {
            snapStatus.textContent = 'Instalado';
            snapStatus.className = 'feature-status available';
        } else {
            snapStatus.textContent = 'Indisponível';
            snapStatus.className = 'feature-status unavailable';
        }
    } catch (err) {
        document.getElementById('distro-name').textContent = 'Erro ao detectar';
    }
}

async function loadUpdateSteps() {
    try {
        const steps = await GetUpdateSteps();
        state.steps = steps;
        renderUpdateOptions(steps);
    } catch (err) {
        console.error('Erro ao carregar opções:', err);
    }
}

function renderUpdateOptions(steps) {
    const container = document.getElementById('update-options');
    container.innerHTML = '';

    if (steps.length === 0) {
        container.innerHTML = '<p style="color:var(--text-muted)">Nenhuma opção de atualização disponível para este sistema.</p>';
        return;
    }

    steps.forEach(step => {
        const el = document.createElement('div');
        el.className = 'update-option';
        el.dataset.stepId = step.id;
        el.innerHTML = `
            <div class="update-checkbox"></div>
            <span class="update-option-label">${step.label}</span>
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
}

async function startUpdate() {
    if (state.running || state.selectedSteps.size === 0) return;

    state.running = true;
    const btn = document.getElementById('btn-update');
    btn.disabled = true;
    btn.innerHTML = '<span class="btn-icon">⏳</span> Atualizando...';

    const progressArea = document.getElementById('progress-area');
    progressArea.classList.remove('hidden');

    const stepsDiv = document.getElementById('progress-steps');
    stepsDiv.innerHTML = '';

    const ids = Array.from(state.selectedSteps);

    ids.forEach(id => {
        const step = state.steps.find(s => s.id === id);
        if (!step) return;

        const el = document.createElement('div');
        el.className = 'progress-step';
        el.id = 'progress-' + id;
        el.innerHTML = `
            <div class="progress-step-header">
                <span class="progress-step-name">${step.label}</span>
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
        addLog('Erro ao iniciar: ' + err, true);
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

    EventsOn('update:complete', () => {
        state.running = false;
        const btn = document.getElementById('btn-update');
        btn.innerHTML = '<span class="btn-icon">▶</span> Iniciar Atualização';
        btn.disabled = state.selectedSteps.size === 0;

        const stepsDiv = document.getElementById('progress-steps');
        const banner = document.createElement('div');
        banner.className = 'complete-banner';
        banner.innerHTML = `
            <span class="complete-banner-icon">✅</span>
            <span class="complete-banner-text">Atualização concluída</span>
        `;
        stepsDiv.appendChild(banner);
    });
}

function addLog(text, isError) {
    const log = document.getElementById('log-output');
    const line = document.createElement('div');
    line.className = 'log-line' + (isError ? ' error' : '');
    line.textContent = text;
    log.appendChild(line);
    log.scrollTop = log.scrollHeight;
}

function confirmAction(action) {
    const labels = {
        reboot: { title: 'Reiniciar Sistema', desc: 'O sistema será reiniciado. Salve seu trabalho antes de continuar.' },
        shutdown: { title: 'Desligar Sistema', desc: 'O sistema será desligado. Salve seu trabalho antes de continuar.' },
    };

    const info = labels[action];
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.innerHTML = `
        <div class="modal">
            <h3>${info.title}</h3>
            <p>${info.desc}</p>
            <div class="modal-actions">
                <button class="btn-cancel" id="modal-cancel">Cancelar</button>
                <button class="btn-danger" id="modal-confirm">Confirmar</button>
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
            addLog('Erro: ' + err, true);
        }
    });
}
