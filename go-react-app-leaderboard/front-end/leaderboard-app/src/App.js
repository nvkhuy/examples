import React, { useState, useEffect, useRef, useCallback } from 'react';
import useFetchData from './hooks/useFetchData';
import { updatePoints } from './hooks/updatePoints';
import LeaderboardTable from "./components/LeaderboardTable";
import Podium from "./components/Podium";

function App() {
    const [refresh, setRefresh] = useState(false);
    const [limit, setLimit] = useState(10);
    const [data, setData] = useState([]);
    const [isLoading, setIsLoading] = useState(false);

    // Fetch data using the custom hook
    const { total: totalAvailableData, data: fetchData } = useFetchData(`http://localhost:8080/users/rankings?limit=${limit}&refresh=${refresh}`);

    useEffect(() => {
        if (fetchData && fetchData.length > 0) {
            setData(fetchData); // Update state with new data
            setIsLoading(false);
        }
    }, [fetchData]);

    // Handler for updating points
    const handleUpdatePoints = (index) => {
        const newPoints = parseInt(prompt('Enter new points:', data[index].points));
        if (!isNaN(newPoints)) {
            const userId = data[index].id;
            updatePoints(userId, newPoints, data, setData);
            setTimeout(() => setRefresh(prev => !prev), 100); // Refresh delay
        }
    };

    // Infinite scroll using IntersectionObserver
    const observer = useRef();
    const lastElementRef = useCallback(node => {
        if (isLoading) return;
        if (observer.current) observer.current.disconnect();
        observer.current = new IntersectionObserver(entries => {
            if (entries[0].isIntersecting && limit < totalAvailableData) {
                setIsLoading(true);
                setLimit(prevLimit => Math.min(prevLimit + 10, totalAvailableData)); // Ensure limit does not exceed total data
            }
        });
        if (node) observer.current.observe(node);
    }, [isLoading, limit, totalAvailableData]);

    return (
        <div className="flex flex-col items-center my-5">
            <h1 className="text-2xl font-bold mb-5">Leaderboard - Top {limit}</h1>
            <LeaderboardTable data={data} handleUpdatePoints={handleUpdatePoints} lastElementRef={lastElementRef} />
            <Podium data={data.slice(0, 3)} />
        </div>
    );
}

export default App;
