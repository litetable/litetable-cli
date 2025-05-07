import AppLayout from "./layout";
import QueryBuilder from '@/components/query-builder';

function App() {
  return (
    <AppLayout>
      <div>
        <h1 className={"text-3xl pb-4"}>Welcome to LiteTable!</h1>
        <QueryBuilder />
        {/*<QueryInput />*/}
      </div>
    </AppLayout>
  );
}

export default App;
