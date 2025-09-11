import { Shield } from "lucide-react"

const Logo = ({ className = "" }) => {
  return (
    <div className={`flex items-center space-x-3 ${className}`}>
      <div className="p-2 bg-indigo-600 rounded-lg">
        <Shield className="w-8 h-8 text-white" />
      </div>
      <span className="text-2xl font-bold text-slate-200">Sentinel</span>
    </div>
  )
}

export default Logo
