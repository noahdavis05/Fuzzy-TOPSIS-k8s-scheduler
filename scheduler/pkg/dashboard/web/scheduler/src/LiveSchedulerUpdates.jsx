
function LiveSchedulerUpdates({ data }) {

    return (
        <>
        <div className="card">
            <div className="card-header">
                <h3>Recently Scheduled pods</h3>
            </div>
            <div className="card-body">
                <pre>{data ? JSON.stringify(data, null, 2) : "No updates yet"}</pre>
            </div>
        </div>
        </>
    )
}

export default LiveSchedulerUpdates
