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
    const fullPortCheckbox = document.getElementById('fullPortScan');

    if (localStorage.getItem('followRedirects') === 'true') {
        followCheckbox.checked = true;
    }
    followCheckbox.addEventListener('change', function () {
        localStorage.setItem('followRedirects', String(followCheckbox.checked));
    });

    if (localStorage.getItem('fullPortScan') === 'true') {
        fullPortCheckbox.checked = true;
    }
    fullPortCheckbox.addEventListener('change', function () {
        localStorage.setItem('fullPortScan', String(fullPortCheckbox.checked));
    });

    // Redirect if someone pastes a /?q=... URL
    const q = (new URLSearchParams(window.location.search).get('q') || '').trim();
    if (q) {
        const dest = new URLSearchParams({ q });
        if (followCheckbox.checked) dest.set('follow', '1');
        if (fullPortCheckbox.checked) dest.set('fullport', '1');
        window.location.replace('/scan?' + dest.toString());
    }
}

function initScanPage() {
    const params = new URLSearchParams(window.location.search);
    const q = (params.get('q') || '').trim();
    const follow = params.get('follow') === '1';
    const fullPort = params.get('fullport') === '1';

    const form = document.getElementById('scanForm');
    const input = document.getElementById('targetInput');
    const button = document.getElementById('scanButton');
    const followCheckbox = document.getElementById('followRedirects');
    const fullPortCheckbox = document.getElementById('fullPortScan');
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
        if (fullPortCheckbox.checked) dest.set('fullport', '1');
        window.location.href = '/scan?' + dest.toString();
    });

    // Auto-start scan from URL params on page load
    if (q) {
        runScan(q, follow, fullPort, {
            button, input, spinner, spinnerText,
            logsSection, logsOutput,
            resultsSection, scanResults,
            errorSection, errorMessage,
        });
    }
}

async function runScan(target, follow, fullPortScan, els) {
    const { button, input, spinner, spinnerText, logsSection, logsOutput,
            resultsSection, scanResults, errorSection, errorMessage } = els;

    input.disabled = true;
    button.disabled = true;
    button.textContent = 'Scanning...';

    try {
        const enqueueResp = await fetch('/api/scan', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ target, follow, full_port_scan: fullPortScan }),
        });
        const enqueueData = await enqueueResp.json();
        if (!enqueueResp.ok) throw new Error(enqueueData.error || enqueueResp.statusText);

        await streamScanJob(enqueueData.job_id, { spinner, spinnerText, logsSection, logsOutput, resultsSection, scanResults });
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

function streamScanJob(jobID, { spinner, spinnerText, logsSection, logsOutput, resultsSection, scanResults }) {
    const timeoutMs = 180000;

    return new Promise((resolve, reject) => {
        const source = new EventSource('/api/scan/' + encodeURIComponent(jobID) + '/events');

        const timeoutTimer = setTimeout(function () {
            source.close();
            reject(new Error('Scan timed out'));
        }, timeoutMs);

        source.onmessage = function (event) {
            const data = JSON.parse(event.data);
            const status = (data.status || '').toLowerCase();

            if (data.stderr) {
                logsSection.style.display = 'block';
                logsOutput.textContent = cleanAnsi(data.stderr);
                logsOutput.scrollTop = logsOutput.scrollHeight;
            }

            if (status === 'running') {
                spinnerText.textContent = getLastLine(data.stderr || '') || 'Scanning...';
            } else if (status === 'completed') {
                clearTimeout(timeoutTimer);
                source.close();
                spinner.style.display = 'none';
                resultsSection.style.display = 'block';
                scanResults.textContent = data.results || '';
                resolve();
            } else if (status === 'failed' || status === 'timeout') {
                clearTimeout(timeoutTimer);
                source.close();
                reject(new Error(data.error || 'Scan ' + status));
            }
        };

        // EventSource retries transient connection drops on its own; the
        // timeout above is the only failure signal needed here.
        source.onerror = function () {};
    });
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
