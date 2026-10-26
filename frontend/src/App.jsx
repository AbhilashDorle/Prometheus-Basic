import { useState } from 'react'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  async function changeCount(action) {
  try {
    console.log(count)
    const response = await fetch(
      `http://localhost:3001/api/${action}`,
      { method: 'POST' }
    )

    if (!response.ok) {
      throw new Error(`Request failed: ${response.status}`)
    }

    const data = await response.json()
    console.log(count) 
    setCount(data.count)
  } catch (error) {
    console.error('Could not update count:', error)
  }
}

  return (
    <main className="counter-page">
      <section className="counter-card" aria-labelledby="counter-title">
        <p className="counter-label">ONE STEP AT A TIME</p>
        <h1 id="counter-title">Make it count.</h1>
        <output className="counter-value" aria-live="polite">{count}</output>
        <div className="counter-controls">
          <button 
            onClick={() => changeCount('decrement')} 
            aria-label="Decrease count">
              −
          </button>
          <button 
            onClick={() => changeCount('increment')}
            aria-label="Increase count">
              +
          </button>
        </div>
      </section>
    </main>
  )
}

export default App