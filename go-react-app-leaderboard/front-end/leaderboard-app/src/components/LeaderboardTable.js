import React from 'react';
import LeaderboardRow from './LeaderboardRow';

function LeaderboardTable({ data, handleUpdatePoints, lastElementRef }) {
    return (
        <div className="w-full max-w-4xl shadow-md relative overflow-y-auto border border-gray-200 rounded" style={{ maxHeight: '30rem' }}>
            <table className="table-auto w-full">
                <thead className="bg-gray-100 sticky top-0 z-10">
                <tr>
                    <th className="px-4 py-2 text-left">Place</th>
                    <th className="px-4 py-2 text-left">Name</th>
                    <th className="px-4 py-2 text-left">Rank Change</th>
                    <th className="px-4 py-2 text-left">Points</th>
                </tr>
                </thead>
                <tbody>
                {data.map((entry, index) => (
                    <LeaderboardRow
                        key={entry.id}
                        entry={entry}
                        index={index}
                        onEdit={() => handleUpdatePoints(index)}
                        lastElementRef={index === data.length - 1 ? lastElementRef : null}
                    />
                ))}
                </tbody>
            </table>
        </div>
    );
}

export default LeaderboardTable;
