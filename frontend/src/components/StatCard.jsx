const StatCard = ({ title, value,   change, changeType = "positive" }) => {
  return (
    <div className="bg-slate-800 rounded-xl p-6 shadow-lg border border-slate-700/50 hover:border-indigo-500/30 transition-all duration-300">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="p-2 bg-indigo-600/20 rounded-lg">
            <Icon className="w-6 h-6 text-indigo-400" />
          </div>
          <div>
            <p className="text-slate-400 text-sm font-medium">{title}</p>
            <p className="text-slate-200 text-2xl font-bold">{value}</p>
          </div>
        </div>
        {change && (
          <div className={`text-sm font-medium ${changeType === "positive" ? "text-green-400" : "text-red-400"}`}>
            {changeType === "positive" ? "+" : ""}
            {change}%
          </div>
        )}
      </div>
    </div>
  )
}

export default StatCard
