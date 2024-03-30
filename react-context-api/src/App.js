import {useContext} from 'react';
import './App.css';
import Button from './Button';
import {ThemeContext} from './store/ThemeContext';

function App() {
    const {isDarkTheme} = useContext(ThemeContext)

    return (
        <div className={`App ${isDarkTheme ? 'darkTheme' : 'lightTheme'}`}>
            <h1>Theme Context API</h1>
            <p>Them Context API example</p>
            <Button/>
        </div>
    );
}

export default App;
