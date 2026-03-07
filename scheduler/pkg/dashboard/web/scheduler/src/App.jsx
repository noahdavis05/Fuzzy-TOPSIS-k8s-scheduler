import 'bootstrap/dist/css/bootstrap.min.css'
import LiveTelemetryGraphs from './LiveTelemetryGraphs'
import LiveTelemetryCache from './LiveTelemetryCache'
import LiveSchedulerUpdates from './LiveSchedulerUpdates'

import { useEffect, useState } from "react";

function App() {

  const [schedulerUpdates, setSchedulerUpdates] = useState(null);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8090/ws");

    ws.onopen = () => {
      console.log("Connected to WebSocket");
    };

    ws.onmessage = (event) => {
      // event.data is plain text
      const parsed = JSON.parse(event.data);
      setSchedulerUpdates(parsed);
    };

    ws.onerror = (err) => console.error("WebSocket error:", err);
    ws.onclose = () => console.log("WebSocket disconnected");

    return () => ws.close();
  }, []);

  return (
    <>
    <div className='nav'>
      Scheduler
    </div>
     <div className='row'>
        <div className='col'>
          <LiveTelemetryGraphs></LiveTelemetryGraphs>
        </div>
        <div className='col'>
          <LiveTelemetryCache></LiveTelemetryCache>
        </div>
        <div className='col'>
          <LiveSchedulerUpdates data={schedulerUpdates}></LiveSchedulerUpdates>
        </div>
     </div>
    </>
  )
}

export default App
