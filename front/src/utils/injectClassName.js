import h from '/utils/h';

const injectClassName = (element, className, properties = {}) => {
  if (!element || !element.props) {
    return element;
  }
  const classNames = [element.props.className, className].filter(Boolean);
  return h(element.type, { ...element.props, ...properties, className: classNames.join(' ') });
};

export default injectClassName;
