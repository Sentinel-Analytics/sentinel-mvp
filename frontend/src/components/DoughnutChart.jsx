import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js"
import { Doughnut } from "react-chartjs-2"

ChartJS.register(ArcElement, Tooltip, Legend)

const DoughnutChart = ({ title, data, labels }) => {
  const colors = [
    "rgba(79, 70, 229, 0.8)",
    "rgba(99, 102, 241, 0.8)",
    "rgba(129, 140, 248, 0.8)",
    "rgba(165, 180, 252, 0.8)",
    "rgba(196, 181, 253, 0.8)",
  ]

  const chartData = {
    labels,
    datasets: [
      {
        data,
        backgroundColor: colors,
        borderColor: colors.map((color) => color.replace("0.8", "1")),
        borderWidth: 2,
        hoverBorderWidth: 3,
      },
    ],
  }

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: "bottom",
        labels: {
          color: "#94a3b8",
          font: {
            size: 12,
          },
          padding: 20,
          usePointStyle: true,
        },
      },
      tooltip: {
        backgroundColor: "rgba(15, 23, 42, 0.95)",
        titleColor: "#e2e8f0",
        bodyColor: "#e2e8f0",
        borderColor: "rgba(79, 70, 229, 0.5)",
        borderWidth: 1,
        cornerRadius: 8,
        displayColors: false,
      },
    },
    cutout: "60%",
  }

  return (
    <div className="bg-slate-800 rounded-xl p-6 shadow-lg border border-slate-700/50 hover:border-indigo-500/30 transition-all duration-300">
      <h3 className="text-slate-200 text-lg font-semibold mb-4">{title}</h3>
      <div className="h-64">
        <Doughnut data={chartData} options={options} />
      </div>
    </div>
  )
}

export default DoughnutChart
