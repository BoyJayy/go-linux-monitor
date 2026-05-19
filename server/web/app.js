let selectedHostID = null;
let devicesCache = [];

let cpuChart;
let memChart;
let netChart;

const refreshBtn = document.getElementById("refreshBtn");
const apiStatus = document.getElementById("apiStatus");

const nodesCount = document.getElementById("nodesCount");
const onlineCount = document.getElementById("onlineCount");
const offlineCount = document.getElementById("offlineCount");
const devicesList = document.getElementById("devicesList");

const selectedHost = document.getElementById("selectedHost");
const selectedState = document.getElementById("selectedState");

const cpuValue = document.getElementById("cpuValue");
const memValue = document.getElementById("memValue");
const rxValue = document.getElementById("rxValue");
const txValue = document.getElementById("txValue");

const cpuMeter = document.getElementById("cpuMeter");
const memMeter = document.getElementById("memMeter");

const historyCount = document.getElementById("historyCount");
const samplesTable = document.getElementById("samplesTable");

refreshBtn.addEventListener("click", () => refreshAll());

initCharts();
refreshAll();
setInterval(refreshAll, 2000);

async function refreshAll() {
    try {
        setAPIStatus("checking", "warn");

        const devices = await fetchJSON("/api/v1/devices");
        devicesCache = devices;

        renderDevices(devices);

        if (!selectedHostID && devices.length > 0) {
            const online = devices.find((device) => device.online);
            selectedHostID = online ? online.host_id : devices[0].host_id;
        }

        if (selectedHostID) {
            await refreshSelectedNode(selectedHostID);
        }

        setAPIStatus("online", "ok");
    } catch (error) {
        console.error(error);
        setAPIStatus("error", "bad");
    }
}

async function refreshSelectedNode(hostID) {
    const latest = await fetchJSON(`/api/v1/devices/latest?host_id=${encodeURIComponent(hostID)}`);
    const history = await fetchJSON(`/api/v1/devices/metrics?host_id=${encodeURIComponent(hostID)}&limit=80`);

    const selectedDevice = devicesCache.find((device) => device.host_id === hostID);
    renderLatest(selectedDevice, latest);
    renderHistory(history);
}

function renderDevices(devices) {
    const online = devices.filter((device) => device.online).length;
    const offline = devices.length - online;

    nodesCount.textContent = `${devices.length} registered`;
    onlineCount.textContent = online;
    offlineCount.textContent = offline;

    if (devices.length === 0) {
        devicesList.innerHTML = `<div class="empty">No devices registered</div>`;
        return;
    }

    devicesList.innerHTML = "";

    for (const device of devices) {
        const item = document.createElement("div");
        item.className = "node-item";

        if (device.host_id === selectedHostID) {
            item.classList.add("active");
        }

        const stateClass = device.online ? "online" : "offline";
        const stateText = device.online ? "online" : "offline";

        item.innerHTML = `
            <div class="node-main">
                <div class="node-host">${escapeHTML(shortHost(device.host_id))}</div>
                <span class="node-state ${stateClass}">${stateText}</span>
            </div>
            <div class="node-meta">
                last_seen=${formatDateTime(device.last_seen_at)}
            </div>
        `;

        item.addEventListener("click", async () => {
            selectedHostID = device.host_id;
            renderDevices(devicesCache);
            await refreshSelectedNode(selectedHostID);
        });

        devicesList.appendChild(item);
    }
}

function renderLatest(device, latest) {
    selectedHost.textContent = latest.host_id || selectedHostID || "No node selected";

    const online = device && device.online;
    selectedState.className = `node-state ${online ? "online" : "offline"}`;
    selectedState.textContent = online ? "online" : "offline";

    const cpu = latest.cpu?.total_pct ?? 0;
    const mem = latest.mem?.used_pct ?? 0;
    const rx = latest.network?.rx_bps_total ?? 0;
    const tx = latest.network?.tx_bps_total ?? 0;

    cpuValue.textContent = `${cpu.toFixed(2)}%`;
    memValue.textContent = `${mem.toFixed(2)}%`;
    rxValue.textContent = formatBps(rx);
    txValue.textContent = formatBps(tx);

    cpuMeter.style.width = `${clamp(cpu, 0, 100)}%`;
    memMeter.style.width = `${clamp(mem, 0, 100)}%`;
}

function renderHistory(history) {
    historyCount.textContent = `${history.length} points`;

    const ordered = [...history].reverse();

    const labels = ordered.map((item) => formatTime(item.timestamp));
    const cpu = ordered.map((item) => item.cpu?.total_pct ?? 0);
    const mem = ordered.map((item) => item.mem?.used_pct ?? 0);
    const rx = ordered.map((item) => bytesToKB(item.network?.rx_bps_total ?? 0));
    const tx = ordered.map((item) => bytesToKB(item.network?.tx_bps_total ?? 0));

    updateChart(cpuChart, labels, [cpu]);
    updateChart(memChart, labels, [mem]);
    updateChart(netChart, labels, [rx, tx]);

    renderSamples(history.slice(0, 8));
}

function renderSamples(samples) {
    if (!samples || samples.length === 0) {
        samplesTable.innerHTML = `<tr><td colspan="5" class="empty-cell">No data</td></tr>`;
        return;
    }

    samplesTable.innerHTML = "";

    for (const sample of samples) {
        const row = document.createElement("tr");

        row.innerHTML = `
            <td>${formatTime(sample.timestamp)}</td>
            <td>${formatPercent(sample.cpu?.total_pct)}</td>
            <td>${formatPercent(sample.mem?.used_pct)}</td>
            <td>${formatBps(sample.network?.rx_bps_total)}</td>
            <td>${formatBps(sample.network?.tx_bps_total)}</td>
        `;

        samplesTable.appendChild(row);
    }
}

function initCharts() {
    cpuChart = createLineChart("cpuChart", ["CPU %"], ["#7393b3"], 0, 100);
    memChart = createLineChart("memChart", ["Memory %"], ["#6aa36f"], 0, 100);
    netChart = createLineChart("netChart", ["RX KB/s", "TX KB/s"], ["#c3a45c", "#b86b6b"], 0, undefined);
}

function createLineChart(canvasID, datasetLabels, colors, min, max) {
    const ctx = document.getElementById(canvasID);

    return new Chart(ctx, {
        type: "line",
        data: {
            labels: [],
            datasets: datasetLabels.map((label, index) => ({
                label,
                data: [],
                borderColor: colors[index],
                backgroundColor: "transparent",
                borderWidth: 1.5,
                pointRadius: 0,
                tension: 0.18,
            })),
        },
        options: {
            animation: false,
            responsive: true,
            maintainAspectRatio: false,
            interaction: {
                intersect: false,
                mode: "index",
            },
            plugins: {
                legend: {
                    labels: {
                        color: "#7f8b96",
                        boxWidth: 10,
                        font: {
                            family: "monospace",
                        },
                    },
                },
            },
            scales: {
                x: {
                    ticks: {
                        color: "#56616c",
                        maxTicksLimit: 8,
                    },
                    grid: {
                        color: "#1d252d",
                    },
                },
                y: {
                    min,
                    max,
                    ticks: {
                        color: "#56616c",
                    },
                    grid: {
                        color: "#1d252d",
                    },
                },
            },
        },
    });
}

function updateChart(chart, labels, datasets) {
    chart.data.labels = labels;

    datasets.forEach((data, index) => {
        chart.data.datasets[index].data = data;
    });

    chart.update();
}

async function fetchJSON(url) {
    const response = await fetch(url);

    if (!response.ok) {
        throw new Error(`${url}: ${response.status}`);
    }

    return response.json();
}

function setAPIStatus(text, cls) {
    apiStatus.textContent = text;
    apiStatus.className = cls;
}

function shortHost(hostID) {
    if (!hostID) return "unknown";
    if (hostID.length <= 34) return hostID;
    return `${hostID.slice(0, 31)}...`;
}

function formatPercent(value) {
    if (typeof value !== "number") return "—";
    return `${value.toFixed(2)}%`;
}

function formatBps(value) {
    if (typeof value !== "number") return "—";

    if (value >= 1024 * 1024) {
        return `${(value / 1024 / 1024).toFixed(2)} MB/s`;
    }

    if (value >= 1024) {
        return `${(value / 1024).toFixed(2)} KB/s`;
    }

    return `${value.toFixed(0)} B/s`;
}

function bytesToKB(value) {
    if (typeof value !== "number") return 0;
    return value / 1024;
}

function formatTime(value) {
    if (!value) return "—";

    const date = new Date(value);
    return date.toLocaleTimeString();
}

function formatDateTime(value) {
    if (!value) return "—";

    const date = new Date(value);
    return date.toLocaleString();
}

function clamp(value, min, max) {
    return Math.max(min, Math.min(max, value));
}

function escapeHTML(value) {
    return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#039;");
}