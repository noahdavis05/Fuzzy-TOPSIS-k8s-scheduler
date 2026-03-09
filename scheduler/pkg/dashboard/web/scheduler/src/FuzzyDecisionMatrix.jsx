import { useState } from "react";

function FuzzyDecisionMatrix({ matrix, title }) {
  const [collapsed, setCollapsed] = useState(true);

  if (!matrix || !matrix.Data) return null;

  const nodes = Object.keys(matrix.Data);
  const criteria = ["CPU", "CPU RANGE", "RAM", "RAM RANGE"];

  return (
    <div style={{ marginBottom: "2rem" }}>
      
      <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
        <h5 style={{ margin: 0 }}>{title}</h5>

        <button className="btn btn-primary"
          onClick={() => setCollapsed(!collapsed)}
          style={{
            padding: "2px 8px",
            cursor: "pointer",
            fontSize: "0.8rem"
          }}
        >
          {collapsed ? "Show" : "Hide"}
        </button>
      </div>

      {!collapsed && (
        <table
          border="1"
          cellPadding="6"
          style={{ borderCollapse: "collapse", width: "100%", marginTop: "10px" }}
        >
          <thead>
            <tr>
              <th>Node</th>
              {criteria.map((c) => (
                <th key={c}>{c}</th>
              ))}
            </tr>
          </thead>

          <tbody>
            {nodes.map((node) => (
              <tr key={node}>
                <td><b>{node}</b></td>

                {criteria.map((criterion) => {
                  const value = matrix.Data[node]?.[criterion];

                  if (!value || value.A === undefined) {
                    return <td key={criterion}>-</td>;
                  }

                  return (
                    <td key={criterion}>
                      ({value.A.toFixed(2)}, {value.B.toFixed(2)}, {value.C.toFixed(2)})
                    </td>
                  );
                })}
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default FuzzyDecisionMatrix;