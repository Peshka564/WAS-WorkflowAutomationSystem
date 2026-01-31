import { useState } from 'react';

export function App() {
  const [val, setVal] = useState('');
  const onClick = async () => {
    const res = await fetch('http://localhost:3000', {
      method: 'GET',
    });
    setVal(await res.text());
  };

  return (
    <>
      <p>{val}</p>
      <button onClick={onClick}>Save workflow</button>
    </>
  );
}
