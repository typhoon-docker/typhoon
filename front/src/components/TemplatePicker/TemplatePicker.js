import h from '/utils/h';

import templates from '/utils/templates.json';

import node from './node.png';
import php from './php.png';
import python3 from './python3.svg';
import raw from './static.svg';
import cra from './create-react-app.svg';
import wordpress from './wordpress.svg';

import { img, label, wrapper, name } from './TemplatePicker.css';

import { Radio } from '/components/Input';

const images = {
  node,
  php,
  python3,
  static: raw,
  'create-react-app': cra,
  wordpress,
};

const TemplatePicker = ({ onSelect }) => {
  return (
    <>
      <h2>En quoi as-tu codé ton projet ?</h2>
      <div className={wrapper}>
        {Object.values(templates).map(template => (
          <div key={template.template_id}>
            <Radio
              id={template.template_id}
              name="template_id"
              value={template.template_id}
              className={label}
              onChange={() => onSelect(template)}
            >
              <>
                <img
                  src={images[template.template_id]}
                  alt={template.template_id}
                  className={img}
                  title={template.template_id}
                />
                <span className={name}>{template.template_id}</span>
              </>
            </Radio>
          </div>
        ))}
      </div>
    </>
  );
};

export default TemplatePicker;
