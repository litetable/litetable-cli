import AppLayout from "./layout";
import QueryBuilder from '@/components/query-builder';

function App() {
  return (
    <AppLayout>
      <div className="container mx-auto max-w-[1200px] py-8">
        <h1 className="text-3xl font-bold text-center mb-6">Welcome to LiteTable!</h1>
        <QueryBuilder />
      </div>
    </AppLayout>
  );
}

export default App;
