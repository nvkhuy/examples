import React, { forwardRef } from 'react';
import getOrdinalSuffix from "../common/getOrdinalSuffix";

const LeaderboardRow = forwardRef(({ entry, index, onEdit, lastElementRef }, ref) => {
    const rankingChangeDisplay = () => {
        if (entry.ranking_change > 0) {
            return `+${entry.ranking_change} RANK`;
        } else if (entry.ranking_change < 0) {
            return `${entry.ranking_change} RANK`;
        } else {
            return "—";
        }
    };

    const rankingChangeColor = entry.ranking_change > 0 ? 'text-green-500' : entry.ranking_change < 0 ? 'text-red-500' : 'text-gray-500';

    return (
        <tr className="border-b last:border-b-0 hover:bg-gray-50" ref={lastElementRef}>
            <td className="px-4 py-2 font-bold">
                <span className={`text-lg font-bold ${entry.ranking_change > 0 ? 'text-green-500' : entry.ranking_change < 0 ? 'text-red-500' : 'text-gray-500'}`}>
                    {getOrdinalSuffix(index + 1)}
                </span>
            </td>
            <td className="px-4 py-2 flex items-center space-x-3">
                <img className="w-10 h-10 rounded-full" src="/default_user.png" alt="avatar" />
                <span>{entry.name}</span>
            </td>
            <td className={`px-4 py-2 font-semibold ${rankingChangeColor}`}>
                {rankingChangeDisplay()}
            </td>
            <td className="px-4 py-2 text-lg font-bold flex items-center">
                <span>{entry.points} PTS</span>
                <button
                    className="ml-3 px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-700"
                    onClick={onEdit}
                >
                    Edit
                </button>
            </td>
        </tr>
    );
});

export default LeaderboardRow;
