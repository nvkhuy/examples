// updatePoints.js

import useFetchData from "./useFetchData";

/**
 * Function to update user points.
 * @param {string} userId - The ID of the user whose points are being updated.
 * @param {number} newPoints - The new points value for the user.
 * @param {Array} data - The current list of user data.
 * @param {Function} setData - The state setter function to update user data.
 */
export const updatePoints = (userId, newPoints, data, setData) => {
    fetch(`http://localhost:8080/users/${userId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            points: newPoints
        }),
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to update points');
            }
            return response.json();
        })
        .then(updatedUser => {
            // Update the local state with the new points
            const updatedData = data.map(user =>
                user.id === userId ? {...user, points: updatedUser.points} : user
            );
            setData(updatedData);
        })
        .catch(error => console.error('Error updating points:', error));
};

export default updatePoints;