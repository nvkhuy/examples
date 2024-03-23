import "./App.css";
import { UsingTheRightWay } from "./errorHanlding/usingTheRightWay";
import { ErrorBoundary } from "react-error-boundary";
import { ProductsFetchingError } from "./components/errors/ErrorBoundrayComponent";

function App() {
  return (
    <div className="App">
      {/* <UsingState /> */}
      {/* <StandardErrorBoundary>
        <UsingErrorBoundaries />
      </StandardErrorBoundary> */}

      <ErrorBoundary
        FallbackComponent={ProductsFetchingError}
        onError={() => console.log("Error happened!")}
      >
        <UsingTheRightWay />
      </ErrorBoundary>
    </div>
  );
}

export default App;
