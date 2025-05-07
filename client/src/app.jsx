import AppLayout from "./layout";
import QueryBuilder from '@/components/query-builder';

function App() {
  return (
    <AppLayout>
      <div>
        <h1 className={"text-3xl py-4"}>Welcome to LiteTable!</h1>
        <QueryBuilder />
      </div>
    </AppLayout>
  );
}

export default App;
