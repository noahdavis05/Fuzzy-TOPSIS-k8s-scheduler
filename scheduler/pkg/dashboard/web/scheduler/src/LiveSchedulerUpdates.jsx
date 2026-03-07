
function LiveSchedulerUpdates({ data, setScheduleIndex }) {

    return (
        <>
        <div className="card">
            <div className="card-header">
                <h3>Recently Scheduled pods</h3>
            </div>
            <div className="card-body">
                {data.length === 0 ? (
                    <p>No updates yet</p>
                ) : (
                    data.map((item, originalIndex) => ({ item, originalIndex }))
                    .reverse()
                    .map(({ item, originalIndex }) => (
                        <div
                            key={originalIndex}
                            className="alert alert-info d-flex justify-content-between align-items-center"
                            >
                            <span>
                                Scheduling notification {originalIndex}: pod scheduled to
                                node {item.payload.nodeName}
                            </span>
                            <button
                                className="btn btn-sm btn-primary"
                                onClick={() => setScheduleIndex(originalIndex)}
                            >
                                Details
                            </button>
                        </div>
                    ))
                )}
            </div>
        </div>
        </>
    )
}

export default LiveSchedulerUpdates
