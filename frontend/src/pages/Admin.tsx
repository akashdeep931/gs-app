import { Link } from 'react-router-dom'

export default function Admin() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-slate-950 px-4">
      <h1 className="text-4xl font-bold text-white">Admin</h1>
      <p className="mt-4 text-slate-400">Coming soon</p>
      <Link to="/" className="mt-8 text-sm text-indigo-400 hover:text-indigo-300">
        &larr; Back to home
      </Link>
    </div>
  )
}
