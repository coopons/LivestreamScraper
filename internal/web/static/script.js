let chartInstance = null;
let chartMode = "time"; 

// Dataset colors for each stream
const colors = [
    'rgba(75, 192, 192, 1)',
    'rgba(255, 99, 132, 1)',
    'rgba(54, 162, 235, 1)',
    'rgba(255, 206, 86, 1)',
    'rgba(153, 102, 255, 1)'
];

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

                const datasets = data.map((streamData, idx) => {
                    const color = colors[idx % colors.length];
                    const points = streamData.snapshots.map(p =>
                        chartMode === "time"
                        ? { x: new Date(p.timestamp), y: p.viewer_count }
                        : {
                            x: new Date(p.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
                            y: p.viewer_count
                        }
                    );

                    const avgViewers = Math.round(points.reduce((a, b) => a + b.y, 0) / points.length);
                    const start = new Date(streamData.snapshots[0].timestamp);
                    const end = new Date(streamData.snapshots[streamData.snapshots.length - 1].timestamp);
                    const durationMs = end - start;
                    const durationMin = Math.round(durationMs / 60000);
                    const hours = Math.floor(durationMin / 60);
                    const minutes = durationMin % 60;
                    const durationFormatted = `${String(hours).padStart(2, ' ')}h:${String(minutes).padStart(2, ' ')}m`;

                    stats.push({
                        label: `Stream ${streamData.stream_id}`,
                        avg: avgViewers,
                        duration: durationFormatted,
                        color: color
                    });

                    return {
                        label: `Stream ${streamData.stream_id}`,
                        data: points,
                        borderColor: color,
                        fill: false,
                        hidden: idx !== 0,
                        tension: 0.1,
                        pointRadius: 2
                    };
                });
                    
                // Reset chart each time modal is displayed
                if (chartInstance) {
                    chartInstance.destroy();
                }

                Chart.register(StreamStatsPlugin);
                const maxViewerCount = Math.max(...datasets.flatMap(ds => ds.data.map(point => point.y)));
                console.log(maxViewerCount)
                const ctx = document.getElementById("modalChart").getContext("2d")
                chartInstance = new Chart(ctx, {
                    type: "line",
                    data: {
                        datasets: datasets,
                    },
                    options: {
                        responsive: true,
                        layout: {
                            margin: {
                                top: 40,
                                right:150
                            }
                        },
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

                                    const annotationIndex = index + 1;
                                    if (chart.data.datasets[annotationIndex]) {
                                        chart.data.datasets[annotationIndex].datalabels.display = !dataset.hidden;
                                    }

                                    chart.update();
                                }
                            },
                            streamStats: {
                                stats: stats
                            }
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
            });
    });
});

// Modal closing
const modal = document.getElementById("streamModal");
const closeModal = document.getElementById("closeModal");

if (closeModal && modal) {
    closeModal.addEventListener("click", () => {
        modal.classList.add("hidden");
    });
    window.addEventListener("click", (e) => {
        if (e.target === modal) {
            modal.classList.add("hidden");
        }
    });
}

const toggleButton = document.getElementById("toggleDate")
if (toggleButton) {
    toggleButton.addEventListener("click", () => {
        chartMode = chartMode === "time" ? "category" : "time";
        const selected = document.querySelector(".stream.selected");
        if (selected) selected.click();
    });
}

function control(action) {
    fetch(`/api/control?action=${action}`)
}

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