import sys
import os
import os.path
from datetime import date

def save_plot(fname, plot_name, plot):
    from matplotlib.backends.backend_agg import FigureCanvasAgg as FigureCanvas
    from matplotlib.figure import Figure
    from matplotlib.dates import DateFormatter

    fig = Figure()
    ax = fig.add_subplot(111)
    ax.xaxis.set_major_formatter(DateFormatter('%H:%M'))
    ax.set_xlabel("Time")
    ax.set_ylabel(plot_name)
    fig.set_figheight(20)
    fig.set_figwidth(30)
    fig.autofmt_xdate()

    handles = []
    labels = []
    for graph in plot:
        x, y = plot[graph]
        handles.append(ax.plot(x, y))
        labels.append(graph)

    fig.legend(handles, labels, 1, shadow=True)

    canvas = FigureCanvas(fig)
    canvas.print_figure(fname, dpi=80)

resdir = "results"
files = os.listdir(resdir)
files.sort()
plot = {} # name -> ([dates], [counts])
for f in files:
    t = int(f[:8])
    d = date(t / 10000, (t % 10000) / 100, t % 100)
    data = dict([(k, int(v)) for k, v in [s.split(": ") for s in open(os.path.join(resdir, f), 'r')]])

    for name, count in data.items():
        x, y = plot.get(name, ([], []))
        x.append(d)
        y.append(count)
        plot[name] = (x, y)

save_plot(sys.argv[1], "LOC", {"Total": plot["total"], "Java": plot[".java"]})
