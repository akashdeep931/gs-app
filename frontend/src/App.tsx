import { Routes, Route } from 'react-router-dom'
import Home from './pages/Home'
import User from './pages/User'
import Admin from './pages/Admin'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/user" element={<User />} />
      <Route path="/admin" element={<Admin />} />
    </Routes>
  )
}
