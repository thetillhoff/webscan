document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('scanForm');
    const input = document.getElementById('targetInput');
    const button = document.getElementById('scanButton');
    const followRedirects = document.getElementById('followRedirects');

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

    input.focus();

    if (localStorage.getItem('followRedirects') === 'true') {
        followRedirects.checked = true;
    }

    followRedirects.addEventListener('change', function () {
        localStorage.setItem('followRedirects', String(followRedirects.checked));
    });

    toggleLogsBtn.addEventListener('click', function () {
        logsExpanded = !logsExpanded;
        logsOutput.style.display = logsExpanded ? 'block' : 'none';
        toggleLogsBtn.textContent = logsExpanded ? 'Hide logs' : 'Show logs';
        if (logsExpanded) {
            logsOutput.scrollTop = logsOutput.scrollHeight;
        }
    });

    form.addEventListener('submit', async function (e) {
        e.preventDefault();

        const target = input.value.trim();
        if (!target) return;

        input.disabled = true;
        button.disabled = true;
        button.textContent = 'Scanning...';

        spinner.style.display = 'flex';
        spinnerText.textContent = 'Scanning...';
        resultsSection.style.display = 'none';
        errorSection.style.display = 'none';
        logsSection.style.display = 'none';
        logsOutput.textContent = '';
        scanResults.textContent = '';

        try {
            const enqueueResp = await fetch('/api/scan', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ target: target, follow: followRedirects.checked })
            });

            const enqueueData = await enqueueResp.json();
            if (!enqueueResp.ok) {
                throw new Error(enqueueData.error || enqueueResp.statusText);
            }

            await pollScanJob(enqueueData.job_id);
        } catch (error) {
            spinner.style.display = 'none';
            errorSection.style.display = 'block';
            errorMessage.textContent = error.message;
        } finally {
            input.disabled = false;
            button.disabled = false;
            button.textContent = 'Scan';
        }
    });

    async function pollScanJob(jobID) {
        const pollIntervalMs = 1000;
        const timeoutMs = 180000;
        const startedAt = Date.now();

        while (true) {
            if (Date.now() - startedAt > timeoutMs) {
                throw new Error('Scan timed out');
            }

            const resp = await fetch('/api/scan/' + encodeURIComponent(jobID));
            const data = await resp.json();
            if (!resp.ok) {
                throw new Error(data.error || resp.statusText);
            }

            const status = (data.status || '').toLowerCase();

            if (data.stderr) {
                logsSection.style.display = 'block';
                logsOutput.textContent = cleanAnsi(data.stderr);
                logsOutput.scrollTop = logsOutput.scrollHeight;
            }

            if (status === 'running') {
                const lastLine = getLastLine(data.stderr || '');
                spinnerText.textContent = lastLine || 'Scanning...';
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
            .replace(/\[[0-9;?]*[ -/]*[@-~]/g, '')
            .replace(/\r/g, '')
            .trim();
    }

    function getLastLine(stderr) {
        const lines = cleanAnsi(stderr).split('\n').filter(l => l.trim());
        return lines.length ? lines[lines.length - 1] : '';
    }
});
