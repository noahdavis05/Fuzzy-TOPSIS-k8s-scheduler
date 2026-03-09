import { useState } from "react";

export function TelemetryCacheTable({ telemetryCache, title = "Telemetry Cache" }) {
  const [collapsed, setCollapsed] = useState(true);

  if (!telemetryCache) return null;

  return (
    <div style={{ marginBottom: "2rem" }}>
      <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
        <h5 style={{ margin: 0 }}>{title}</h5>
        <button className="btn btn-primary"
          onClick={() => setCollapsed(!collapsed)}
          style={{ padding: "2px 8px", cursor: "pointer", fontSize: "0.8rem" }}
        >
          {collapsed ? "Show" : "Hide"}
        </button>
      </div>

      {!collapsed && (
        <table border="1" cellPadding="6" style={{ borderCollapse: "collapse", width: "100%", marginTop: "10px" }}>
          <thead>
            <tr>
              <th>Node</th>
              <th>CPU (Low / Mean / High)</th>
              <th>RAM (Low / Mean / High)</th>
              <th>LastScheduled</th>
            </tr>
          </thead>
          <tbody>
            {Object.entries(telemetryCache).map(([node, info]) => (
              <tr key={node}>
                <td><b>{node}</b></td>
                <td>
                  {info.CPU
                    ? `${info.CPU.Low?.toFixed(2)} / ${info.CPU.Mean?.toFixed(2)} / ${info.CPU.High?.toFixed(2)}`
                    : "-"}
                </td>
                <td>
                  {info.RAM
                    ? `${info.RAM.Low?.toFixed(2)} / ${info.RAM.Mean?.toFixed(2)} / ${info.RAM.High?.toFixed(2)}`
                    : "-"}
                </td>
                <td>{info.LastScheduled || "-"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}