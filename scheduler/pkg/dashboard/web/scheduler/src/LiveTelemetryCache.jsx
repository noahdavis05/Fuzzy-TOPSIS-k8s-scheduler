import { useState, useEffect } from 'react';
import { ObjectInspector } from 'react-inspector';

function LiveTelemetryCache({ data }) {
  // State to hold the payload we display
  const [payloadRef, setPayloadRef] = useState(null);

  // Whenever data updates, update payloadRef with the latest payload
  useEffect(() => {
    if (data.length > 0) {
      const latest = data[data.length - 1];

      // Only update if the object reference actually changed
      if (latest.payload !== payloadRef) {
        setPayloadRef(latest.payload);
      }
    }
  }, [data, payloadRef]);

  if (!payloadRef) {
    return (
      <div className="card">
        <div className="card-header">
          <h3>Current Telemetry Cache</h3>
        </div>
        <div className="card-body">
          <p>No telemetry data yet</p>
        </div>
      </div>
    );
  }

  const { unfilteredCache, filteredCache } = payloadRef;

  return (
    <div className="card">
      <div className="card-header">
        <h3>Current Telemetry Cache</h3>
      </div>
      <div className="card-body">

        <div style={{ marginBottom: '1rem' }}>
          <h5>Unfiltered Cache</h5>
          <ObjectInspector
            key="unfiltered"
            data={unfilteredCache}
            expandLevel={1} 
          />
        </div>


        <div>
          <h5>Filtered Cache</h5>
          <ObjectInspector
            key="filtered"
            data={filteredCache}
            expandLevel={1} 
          />
        </div>
      </div>
    </div>
  );
}

export default LiveTelemetryCache;