import { useState, useEffect } from 'react';

/**
 * Custom hook to fetch data from a URL.
 * @param {string} url - The URL to fetch data from.
 * @returns {{ total: number, data: User[] }} - Object containing the total number of users and an array of User objects.
 */
function useFetchData(url) {
    const [data, setData] = useState({ total: 0, data: [] });

    useEffect(() => {
        let isMounted = true; // Track if the component is still mounted

        fetch(url)
            .then((response) => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then((resp) => {
                if (isMounted) {
                    setData({ total: resp.total, data: resp.data }); // Update state with total and data
                }
            })
            .catch((error) => console.error('Error fetching data:', error));

        return () => {
            isMounted = false; // Cleanup to avoid updating state if unmounted
        };
    }, [url]); // Dependency on 'url' ensures the effect reruns only if the url changes

    return data;
}

export default useFetchData;
