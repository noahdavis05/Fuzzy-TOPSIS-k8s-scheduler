import { useState } from "react";

export function NodeScoresTable({ nodeScores, title = "Node Scores" }) {
  const [collapsed, setCollapsed] = useState(true);

  if (!nodeScores) return null;

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
          <thead>
            <tr><th>Node</th><th>Score</th></tr>
          </thead>
          <tbody>
            {Object.entries(nodeScores).map(([node, score]) => (
              <tr key={node}>
                <td>{node}</td>
                <td>{score.toFixed(4)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}