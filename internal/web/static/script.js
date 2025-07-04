let chartInstance = null;

document.querySelectorAll(".stream").forEach(btn => {
    btn.addEventListener("click", () => {
        const streamId = btn.getAttribute("data-stream-id")
        const streamer = btn.getAttribute("data-streamer")

        document.getElementById("modalTitle").textContent = `Viewer Count: ${streamer}`;
        document.getElementById("streamModal").classList.remove("hidden");
        
        fetch(`/api/snapshots?stream_id=${streamId}`)
            .then(res => res.json())
            .then(data => {
                const timestamps = data.map(p =>
                    new Date(p.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
                );
                const counts = data.map(p => p.viewer_count);

                if (chartInstance) {
                    chartInstance.destroy();
                }
            
                const ctx = document.getElementById("modalChart").getContext("2d")
                chartInstance = new Chart(ctx, {
                    type: "line",
                    data: {
                        labels: timestamps,
                        datasets: [{
                        label: `${streamer}'s Logged Viewer Count`,
                        data: counts,
                        borderColor: "rgba(75, 192, 192, 1)",
                        fill: false
                        }]
                    },
                    options: {
                        responsive: true,
                        scales: {
                            y: {
                                beginAtZero: false,
                                title: {
                                    display: true,
                                    text: "Viewers"
                                }
                            },
                            x: {
                                title: {
                                    display: true,
                                    text: "Time"
                                }
                            }
                        }
                    }
                });
            });
    });
});

// Modal closing
const modal = document.getElementById("streamModal");
const closeModal = document.getElementById("closeModal");

closeModal.addEventListener("click", () => {
    modal.classList.add("hidden");
});
window.addEventListener("click", (e) => {
    if (e.target === modal) {
        modal.classList.add("hidden");
    }
});

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