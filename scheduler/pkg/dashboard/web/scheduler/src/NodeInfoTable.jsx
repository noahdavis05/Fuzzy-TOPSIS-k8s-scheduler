import { useState } from "react";

export function NodeInfoTable({ nodeName, cpuRequest, ramRequest, title = "Node Info" }) {
  const [collapsed, setCollapsed] = useState(false);

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
        <table border="1" cellPadding="6" style={{ borderCollapse: "collapse", width: "50%", marginTop: "10px" }}>
          <tbody>
            <tr><td>Node Name</td><td>{nodeName}</td></tr>
            <tr><td>CPU Request</td><td>{cpuRequest}</td></tr>
            <tr><td>RAM Request</td><td>{ramRequest}</td></tr>
          </tbody>
        </table>
      )}
    </div>
  );
}