fetch("/api/snapshots")
  .then(res => res.json())
  .then(data => {
    const streamId = data[0].stream_id;
    const filtered = data.filter(p => p.stream_id === streamId);

    const timestamps = filtered.map(p =>
      new Date(p.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    );
    const counts = filtered.map(p => p.viewer_count);

    const ctx = document.getElementById("viewersChart").getContext("2d");
    new Chart(ctx, {
      type: "line",
      data: {
        labels: timestamps,
        datasets: [{
          label: `Viewer Count for Stream ${streamId}`,
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
