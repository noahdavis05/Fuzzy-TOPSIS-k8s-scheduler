import 'bootstrap/dist/css/bootstrap.min.css'
import LiveTelemetryGraphs from './LiveTelemetryGraphs'
import LiveTelemetryCache from './LiveTelemetryCache'
import LiveSchedulerUpdates from './LiveSchedulerUpdates'
import DetailedScheduleInfo from './DetailedScheduleInfo'

import { useEffect, useState } from "react";

function App() {

  const [scheduleIndex, setScheduleIndex] = useState(null);

  const [schedulerUpdates, setSchedulerUpdates] = useState([]);
  const [telemetryCache, setTelemetryCache] = useState([]);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8090/ws");

    ws.onopen = () => {
      console.log("Connected to WebSocket");
    };

    ws.onmessage = (event) => {
      // event.data is plain text
      let message;
      try {
        message = JSON.parse(event.data); 
      } catch {
        console.log("Message could not be parsed");
        return;
      }

      if (message.subject == "scheduling_update"){
        setSchedulerUpdates(prev => {
          const newArray = [...prev, message];
          return newArray.slice(-500);
        });
      } else if (message.subject == "telemetry_cache"){
        setTelemetryCache(prev => {
          const newArray = [...prev, message];
          return newArray.slice(-500);
        });
      } else if (message.subject == "live_telemetry"){

      }

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
          <LiveTelemetryCache data={telemetryCache}></LiveTelemetryCache>
        </div>
        <div className='col'>
          <LiveSchedulerUpdates data={schedulerUpdates} setScheduleIndex={setScheduleIndex}></LiveSchedulerUpdates>
        </div>
     </div>
     <div className='row'>
        <DetailedScheduleInfo data={schedulerUpdates} index={scheduleIndex}></DetailedScheduleInfo>
     </div>
    </>
  )
}

export default App
