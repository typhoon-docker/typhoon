const cx = (...args) => {
  const className = [];
  args.forEach(arg => {
    if (Array.isArray(arg)) {
      arg.forEach(e => className.push(cx(e)));
    } else if (arg) {
      className.push(arg);
    }
  });
  return className.join(' ');
};

export default cx;
