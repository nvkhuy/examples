import { useContext } from 'react';
import './App.css';
import Button from './Button';
import { ThemeContext } from './store/ThemeContext';

function App() {
  const { isDarkTheme } = useContext(ThemeContext)

  return (
      <div className={`App ${isDarkTheme ? 'darkTheme' : 'lightTheme'}`}>
        <h1>Theme Context Api</h1>
        <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Dolorem quam quisquam enim minus, consectetur dignissimos vero beatae possimus reprehenderit sed officia eveniet obcaecati neque architecto ut, magnam odit optio veniam.</p>
        <Button />
      </div>
  );
}

export default App;
