let selectedHostID = null;

const devicesList = document.getElementById("devicesList");
const selectedDevice = document.getElementById("selectedDevice");

const cpuValue = document.getElementById("cpuValue");
const memoryValue = document.getElementById("memoryValue");
const rxValue = document.getElementById("rxValue");
const txValue = document.getElementById("txValue");

const historyList = document.getElementById("historyList");
const refreshBtn = document.getElementById("refreshBtn");

refreshBtn.addEventListener("click", async () => {
    await loadDevices();

    if (selectedHostID) {
        await loadLatest(selectedHostID);
        await loadHistory(selectedHostID);
    }
});

async function loadDevices() {
    try {
        const res = await fetch("/api/v1/devices");

        if (!res.ok) {
            throw new Error(`failed to load devices: ${res.status}`);
        }

        const devices = await res.json();
        renderDevices(devices);
    } catch (err) {
        devicesList.innerHTML = `<div class="error">${err.message}</div>`;
    }
}

function renderDevices(devices) {
    if (!devices || devices.length === 0) {
        devicesList.textContent = "No devices found";
        return;
    }

    devicesList.innerHTML = "";

    for (const device of devices) {
        const item = document.createElement("div");
        item.className = "device-item";

        if (device.host_id === selectedHostID) {
            item.classList.add("active");
        }

        const statusClass = device.online ? "online" : "offline";
        const statusText = device.online ? "online" : "offline";

        item.innerHTML = `
            <div class="device-host">${escapeHTML(device.host_id)}</div>
            <div class="device-meta">last seen: ${formatTime(device.last_seen_at)}</div>
            <span class="status ${statusClass}">${statusText}</span>
        `;

        item.addEventListener("click", async () => {
            selectedHostID = device.host_id;
            selectedDevice.textContent = `Selected: ${device.host_id}`;

            renderDevices(devices);
            await loadLatest(device.host_id);
            await loadHistory(device.host_id);
        });

        devicesList.appendChild(item);
    }
}

async function loadLatest(hostID) {
    try {
        const res = await fetch(`/api/v1/devices/latest?host_id=${encodeURIComponent(hostID)}`);

        if (!res.ok) {
            throw new Error(`failed to load latest metrics: ${res.status}`);
        }

        const metrics = await res.json();

        cpuValue.textContent = formatPercent(metrics.cpu?.total_pct);
        memoryValue.textContent = formatPercent(metrics.mem?.used_pct);
        rxValue.textContent = formatBps(metrics.network?.rx_bps_total);
        txValue.textContent = formatBps(metrics.network?.tx_bps_total);
    } catch (err) {
        selectedDevice.innerHTML = `<span class="error">${err.message}</span>`;
    }
}

async function loadHistory(hostID) {
    try {
        const res = await fetch(`/api/v1/devices/metrics?host_id=${encodeURIComponent(hostID)}&limit=30`);

        if (!res.ok) {
            throw new Error(`failed to load history: ${res.status}`);
        }

        const history = await res.json();
        renderHistory(history);
    } catch (err) {
        historyList.innerHTML = `<div class="error">${err.message}</div>`;
    }
}

function renderHistory(history) {
    if (!history || history.length === 0) {
        historyList.textContent = "No history found";
        return;
    }

    historyList.innerHTML = "";

    for (const item of history) {
        const row = document.createElement("div");
        row.className = "history-row";

        row.innerHTML = `
            <span>${formatTime(item.timestamp)}</span>
            <span>CPU: ${formatPercent(item.cpu?.total_pct)}</span>
            <span>RAM: ${formatPercent(item.mem?.used_pct)}</span>
            <span>RX: ${formatBps(item.network?.rx_bps_total)}</span>
            <span>TX: ${formatBps(item.network?.tx_bps_total)}</span>
        `;

        historyList.appendChild(row);
    }
}

function formatPercent(value) {
    if (typeof value !== "number") {
        return "—";
    }

    return `${value.toFixed(2)}%`;
}

function formatBps(value) {
    if (typeof value !== "number") {
        return "—";
    }

    if (value >= 1024 * 1024) {
        return `${(value / 1024 / 1024).toFixed(2)} MB/s`;
    }

    if (value >= 1024) {
        return `${(value / 1024).toFixed(2)} KB/s`;
    }

    return `${value.toFixed(0)} B/s`;
}

function formatTime(value) {
    if (!value) {
        return "—";
    }

    const date = new Date(value);
    return date.toLocaleString();
}

function escapeHTML(value) {
    return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#039;");
}

loadDevices();