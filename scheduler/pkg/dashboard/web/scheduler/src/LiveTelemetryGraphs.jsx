import React, { useMemo } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

function LiveTelemetryGraphs({ data }) {
  // pre chart data
  const { nodeNames, labels, cpuSeries, ramSeries } = useMemo(() => {
    const now = new Date();
    const fiveMinutesAgo = new Date(now.getTime() - 5 * 60 * 1000);

    // filter to only use last 5 mins of data
    const filtered = data.filter(d => new Date(d.payload.Timestamp) >= fiveMinutesAgo);
    const labels = filtered.map(d => new Date(d.payload.Timestamp).toLocaleTimeString());

    const nodeNames = filtered.length > 0 ? Object.keys(filtered[0].payload.Data) : [];

    const cpuSeries = {};
    const ramSeries = {};

    nodeNames.forEach(node => {
      cpuSeries[node] = filtered.map(d => d.payload.Data[node]?.CPU ?? 0);
      ramSeries[node] = filtered.map(d => d.payload.Data[node]?.RAM ?? 0);
    });

    return { nodeNames, labels, cpuSeries, ramSeries };
  }, [data]);

  return (
    <div className='card'>
        <div className='card-header'>
            <h2>Live Node Telemetry</h2>
        </div>
        <div className='card-body'>
            {nodeNames.map(node => (
                <div key={node} className="card mb-3">
                    <div className="card-header">
                        <h5>{node} - CPU Usage (%)</h5>
                    </div>
                    <div>
                        <Line
                        data={{
                            labels,
                            datasets: [
                            {
                                label: 'CPU',
                                data: cpuSeries[node],
                                fill: false,
                                borderColor: '#36A2EB',
                                tension: 0.1
                            }
                            ]
                        }}
                        options={{
                            responsive: true,
                            plugins: { legend: { display: false } },
                            scales: { y: { min: 0, max: 100 } }
                        }}
                        />
                    </div>

                    <div className="card-header mt-3">
                        <h5>{node} - RAM Usage (%)</h5>
                    </div>
                    <div>
                        <Line
                        data={{
                            labels,
                            datasets: [
                            {
                                label: 'RAM',
                                data: ramSeries[node],
                                fill: false,
                                borderColor: '#FF6384',
                                tension: 0.1
                            }
                            ]
                        }}
                        options={{
                            responsive: true,
                            plugins: { legend: { display: false } },
                            scales: { y: { min: 0, max: 100 } }
                        }}
                        />
                    </div>
                </div>
            ))}
        </div>
    </div>
  );
}

export default LiveTelemetryGraphs;