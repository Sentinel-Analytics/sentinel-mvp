export function Logo({ className }) {
  return (
    <div className={`flex items-center ${className}`}>
      <div className="relative">
        <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center glow-blue">
          <div className="w-4 h-4 bg-primary-foreground rounded-sm transform rotate-45"></div>
        </div>
        <div className="absolute -top-1 -right-1 w-3 h-3 bg-accent rounded-full animate-pulse"></div>
      </div>
    </div>
  )
}
