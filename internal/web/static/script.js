let chartInstance = null;
let chartMode = "time"; 
let lastSeenRun = null;

// Dataset colors for each stream
const colors = [
    'rgba(75, 192, 192, 1)',
    'rgba(255, 99, 132, 1)',
    'rgba(54, 162, 235, 1)',
    'rgba(255, 206, 86, 1)',
    'rgba(153, 102, 255, 1)'
];

// Plugin to display stream stats below the chart
const StreamStatsPlugin = {
  id: 'streamStats',
  afterDraw(chart, args, pluginOptions) {
    const stats = pluginOptions?.stats;
    if (!stats || !Array.isArray(stats)) return;

    const { ctx, chartArea: { top, right } } = chart;

    ctx.save();
    ctx.font = 'bold 12px sans-serif';

    let yOffset = 20;
    stats.forEach((stat, i) => {
      ctx.fillStyle = stat.color;
      ctx.fillText(`Avg ${stat.avg}👁, ${stat.duration}`, right - 150, top + yOffset);
      yOffset += 18;
    });

    ctx.restore();
  }
};

// Loading modals with with viewer count charts
document.querySelectorAll(".stream").forEach(btn => {
    btn.addEventListener("click", () => {
        const streamId = btn.getAttribute("data-stream-id")
        const streamer = btn.getAttribute("data-streamer")

        document.getElementById("modalTitle").textContent = `Viewer Count: ${streamer}`;
        document.getElementById("streamModal").classList.remove("hidden");

        document.querySelectorAll(".stream").forEach(b => b.classList.remove("selected"));
        btn.classList.add("selected");

        fetch(`/api/snapshots?stream_id=${streamId}`)
            .then(res => res.json())
            .then(data => {

                if (!data || data.length === 0) return;
                
                const stats = [];
                const showOnlyMostRecent = !modalOpen;
                const visibleStreams = window.lastVisibleStreams;

                const datasets = data.map((streamData, idx) => {
                    const points = streamData.snapshots.map(p =>
                        chartMode === "time"
                        ? { x: new Date(p.timestamp), y: p.viewer_count }
                        : {
                            x: new Date(p.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
                            y: p.viewer_count
                        }
                    );
                    
                    const color = colors[idx % colors.length];
                    const avgViewers = Math.round(points.reduce((a, b) => a + b.y, 0) / points.length);
                    const start = new Date(streamData.snapshots[0].timestamp);
                    const end = new Date(streamData.snapshots[streamData.snapshots.length - 1].timestamp);
                    const durationMs = end - start;
                    const durationMin = Math.round(durationMs / 60000);
                    const hours = Math.floor(durationMin / 60);
                    const minutes = durationMin % 60;
                    const durationFormatted = `${String(hours).padStart(2, ' ')}h:${String(minutes).padStart(2, ' ')}m`;
                    const hidden = 
                        visibleStreams !== null && visibleStreams !== undefined
                            ? !visibleStreams.has(streamData.stream_id)
                            : idx !== 0;

                    return {
                        label: `Stream ${streamData.stream_id}`,
                        data: points,
                        borderColor: color,
                        fill: false,
                        hidden: hidden,
                        tension: 0.1,
                        pointRadius: 2,
                        _streamStats: {
                            avg: avgViewers,
                            duration: durationFormatted,
                            color: color
                        }
                    };
                });
                    
                // Reset chart each time modal is displayed
                if (chartInstance) {
                    chartInstance.destroy();
                }

                Chart.register(StreamStatsPlugin);
                const maxViewerCount = Math.max(...datasets.flatMap(ds => ds.data.map(point => point.y)));
                const pluginStats = datasets
                    .filter(ds => !ds.hidden && ds._streamStats)
                    .map(ds => ds._streamStats);

                const ctx = document.getElementById("modalChart").getContext("2d")
                chartInstance = new Chart(ctx, {
                    type: "line",
                    data: {
                        datasets: datasets,
                    },
                    options: {
                        responsive: true,
                        interaction: {
                            mode: `nearest`,
                            axis: `x`,
                            intersect: false
                        },
                        plugins: {
                            legend: {
                                display: true,
                                labels: {
                                    usePointStyle: true
                                },
                                onClick: (e, legendItem, legend) => {
                                    const chart = legend.chart;
                                    const index = legendItem.datasetIndex;
                                    const dataset = chart.data.datasets[index];

                                    dataset.hidden = !dataset.hidden;

                                    // Update stats to only visible ones
                                    chart.options.plugins.streamStats.stats = chart.data.datasets
                                        .filter(ds => !ds.hidden && ds._streamStats)
                                        .map(ds => ds._streamStats);

                                    chart.update();
                                }
                            },
                            streamStats: {
                                stats: pluginStats
                            },
                        },
                        scales: {
                            y: {
                                beginAtZero: false,
                                max: maxViewerCount * 1.2,
                                title: {
                                    display: true,
                                    text: "Viewers"
                                }
                            },
                            x: chartMode === "time" 
                            ? {
                                type: 'time',
                                time: {
                                    unit: "hour",
                                    tooltipFormat: 'MMM d, h:mm a',
                                },
                                title: {
                                    display: true,
                                    text: "Timestamp"
                                }
                            }
                            : {
                                type: "category",
                                title: {
                                    display: true,
                                    text: "Time of Day"
                                }
                            }
                        },
                    },
                    plugins: [StreamStatsPlugin], // enable plugin
                    streamStats: {
                        stats: stats
                    }
                });
                modalOpen = true;
            });
    });
});

// Modal control
const modal = document.getElementById("streamModal");
const closeModal = document.getElementById("closeModal");

if (closeModal && modal) {
    closeModal.addEventListener("click", () => {
        modal.classList.add("hidden");
        window.lastVisibleStreams = null;
    });
    window.addEventListener("click", (e) => {
        if (e.target === modal) {
            modal.classList.add("hidden");
            window.lastVisibleStreams = null;
        }
    });
    modalOpen = false;
}

// Toggle between time and category mode
// This will also update the chart with the last visible streams
const toggleButton = document.getElementById("toggleDate")
if (toggleButton) {
    toggleButton.addEventListener("click", () => {
        if (chartInstance) {
            const visibleStreams = new Set();
            chartInstance.data.datasets.forEach(ds => {
                if (!ds.hidden && ds.label.startsWith("Stream")) {
                    const streamId = ds.label.split(" ")[1]; // "Stream xyz" -> "xyz"
                    visibleStreams.add(streamId);
                }
            });
            // Store in global for reuse with fetch()
            window.lastVisibleStreams = visibleStreams;
        }
        
        chartMode = chartMode === "time" ? "category" : "time";

        const selected = document.querySelector(".stream.selected");
        if (selected) selected.click();
    });
}

// Control collector actions
// This will send a request to the API to start/stop the collector
function control(action) {
    fetch(`/api/control?action=${action}`)
}

// Update countdown timer for next run
// This will fetch the next run time from the API and update the countdown display
function updateCountdown() {
    fetch(`/api/next-run`)
    .then(res => res.json())
    .then(data => {
        if (!data.running) {
            document.getElementById("countdown").textContent = "⏹ Collector is stopped.";
            return;
        }

        const next = new Date(data.next_run);
        const now = new Date();
        const diff = next - now;

        if (diff > 100) { // 100 ms buffer
            const mins = Math.floor(diff / 60000);
            const secs = Math.floor((diff % 60000) / 1000);
            document.getElementById("countdown").textContent = `${mins}m ${secs}s`;
        } else {
            document.getElementById("countdown").textContent = "⏳ Collecting now...";
        }
    });
}

setInterval(updateCountdown, 1000);
updateCountdown();