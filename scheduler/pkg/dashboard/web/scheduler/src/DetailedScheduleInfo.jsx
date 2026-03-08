import { ObjectInspector } from 'react-inspector';


function DetailedScheduleInfo({ data, index}) {

    return (
    <div className="card">
      <div className="card-header">
        <h2>Detailed Schedule Info</h2>
      </div>
      <div className="card-body">
        {index == null || !data ? (
          <p>No Schedule Selected</p>
        ) : (
          <ObjectInspector
            data={data[index].payload}
            expandLevel={2} 
            theme="chromeDark"
            style={{ fontSize: '1.1rem', lineHeight: '1.5' }}
          />
        )}
      </div>
    </div>
  );
}

export default DetailedScheduleInfo