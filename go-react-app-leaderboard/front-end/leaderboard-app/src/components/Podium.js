import React from 'react';
import getOrdinalSuffix from "../common/getOrdinalSuffix";

function Podium({data}) {
    // Sort data by points descending (assuming higher points are better)
    const sortedData = data.sort((a, b) => b.points - a.points).slice(0, 3);

    // Swap the first and second podium places if there are more than one entry
    if (sortedData.length > 1) {
        [sortedData[0], sortedData[1]] = [sortedData[1], sortedData[0]];
    }

    // Style based on position
    const podiumStyles = {
        0: "bg-yellow-100 h-96", // Highest points (center, 1st place)
        1: "bg-gray-100 h-64",  // Second highest points (left, 2nd place)
        2: "bg-red-100 h-48"    // Third highest points (right, 3rd place)
    };

    return (
        <div
            className="flex flex-col md:flex-row justify-center items-end space-y-4 md:space-y-0 md:space-x-4 my-10 bg-white p-6 rounded-lg shadow-md w-full max-w-4xl mx-auto">
            {sortedData.map((user, index) => {
                // Determine position for podium placement
                const position = index === 1 ? 0 : (index === 0 ? 1 : 2);
                return (
                    <div key={user.id}
                         className={`flex flex-col items-center ${podiumStyles[position]} p-4 rounded-lg w-full md:w-1/3`}>
                        <img className="w-16 h-16 rounded-full mb-2" src={user.image || '/default_user.png'}
                             alt={user.name}/>
                        <div className="text-center">
                            <p className="font-semibold">{user.name}</p>
                            <p className="text-gray-600">{getOrdinalSuffix(position + 1)}</p>
                            <p>{`${user.points} PTS - ${user.reward}`}</p>
                        </div>
                    </div>
                );
            })}
        </div>
    );
}

export default Podium;
