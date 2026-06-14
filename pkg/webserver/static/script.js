document.addEventListener('DOMContentLoaded', function () {
    const path = window.location.pathname;
    if (path === '/') {
        initLandingPage();
    } else if (path === '/scan') {
        initScanPage();
    }
});

function initLandingPage() {
    const followCheckbox = document.getElementById('followRedirects');

    if (localStorage.getItem('followRedirects') === 'true') {
        followCheckbox.checked = true;
    }
    followCheckbox.addEventListener('change', function () {
        localStorage.setItem('followRedirects', String(followCheckbox.checked));
    });

    // Redirect if someone pastes a /?q=... URL
    const q = (new URLSearchParams(window.location.search).get('q') || '').trim();
    if (q) {
        const dest = new URLSearchParams({ q });
        if (followCheckbox.checked) dest.set('follow', '1');
        window.location.replace('/scan?' + dest.toString());
    }
}

function initScanPage() {
    const params = new URLSearchParams(window.location.search);
    const q = (params.get('q') || '').trim();
    const follow = params.get('follow') === '1';

    const form = document.getElementById('scanForm');
    const input = document.getElementById('targetInput');
    const button = document.getElementById('scanButton');
    const followCheckbox = document.getElementById('followRedirects');
    const spinner = document.getElementById('spinner');
    const spinnerText = document.getElementById('spinnerText');
    const logsSection = document.getElementById('logsSection');
    const toggleLogsBtn = document.getElementById('toggleLogs');
    const logsOutput = document.getElementById('logsOutput');
    const resultsSection = document.getElementById('resultsSection');
    const scanResults = document.getElementById('scanResults');
    const errorSection = document.getElementById('errorSection');
    const errorMessage = document.getElementById('errorMessage');

    let logsExpanded = false;

    toggleLogsBtn.addEventListener('click', function () {
        logsExpanded = !logsExpanded;
        logsOutput.style.display = logsExpanded ? 'block' : 'none';
        toggleLogsBtn.textContent = logsExpanded ? 'Hide logs' : 'Show logs';
        if (logsExpanded) logsOutput.scrollTop = logsOutput.scrollHeight;
    });

    // Form submit navigates to /scan?q=... — triggers a fresh page load + scan
    form.addEventListener('submit', function (e) {
        e.preventDefault();
        const newQ = input.value.trim();
        if (!newQ) return;
        const dest = new URLSearchParams({ q: newQ });
        if (followCheckbox.checked) dest.set('follow', '1');
        window.location.href = '/scan?' + dest.toString();
    });

    // Auto-start scan from URL params on page load
    if (q) {
        runScan(q, follow, {
            button, input, spinner, spinnerText,
            logsSection, logsOutput,
            resultsSection, scanResults,
            errorSection, errorMessage,
        });
    }
}

async function runScan(target, follow, els) {
    const { button, input, spinner, spinnerText, logsSection, logsOutput,
            resultsSection, scanResults, errorSection, errorMessage } = els;

    input.disabled = true;
    button.disabled = true;
    button.textContent = 'Scanning...';

    try {
        const enqueueResp = await fetch('/api/scan', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ target, follow }),
        });
        const enqueueData = await enqueueResp.json();
        if (!enqueueResp.ok) throw new Error(enqueueData.error || enqueueResp.statusText);

        await pollScanJob(enqueueData.job_id, { spinner, spinnerText, logsSection, logsOutput, resultsSection, scanResults });
    } catch (err) {
        spinner.style.display = 'none';
        errorSection.style.display = 'block';
        errorMessage.textContent = err.message;
    } finally {
        input.disabled = false;
        button.disabled = false;
        button.textContent = 'Scan';
    }
}

async function pollScanJob(jobID, { spinner, spinnerText, logsSection, logsOutput, resultsSection, scanResults }) {
    const pollIntervalMs = 1000;
    const timeoutMs = 180000;
    const startedAt = Date.now();

    while (true) {
        if (Date.now() - startedAt > timeoutMs) {
            throw new Error('Scan timed out');
        }

        const resp = await fetch('/api/scan/' + encodeURIComponent(jobID));
        const data = await resp.json();
        if (!resp.ok) throw new Error(data.error || resp.statusText);

        const status = (data.status || '').toLowerCase();

        if (data.stderr) {
            logsSection.style.display = 'block';
            logsOutput.textContent = cleanAnsi(data.stderr);
            logsOutput.scrollTop = logsOutput.scrollHeight;
        }

        if (status === 'running') {
            spinnerText.textContent = getLastLine(data.stderr || '') || 'Scanning...';
        } else if (status === 'completed') {
            spinner.style.display = 'none';
            resultsSection.style.display = 'block';
            scanResults.textContent = data.results || '';
            return;
        } else if (status === 'failed' || status === 'timeout') {
            throw new Error(data.error || 'Scan ' + status);
        }

        await new Promise(r => setTimeout(r, pollIntervalMs));
    }
}

function cleanAnsi(text) {
    return String(text || '')
        .replace(/\x1b\[[0-9;?]*[ -/]*[@-~]/g, '')
        .replace(/\r/g, '')
        .trim();
}

function getLastLine(stderr) {
    const lines = cleanAnsi(stderr).split('\n').filter(l => l.trim());
    return lines.length ? lines[lines.length - 1] : '';
}
