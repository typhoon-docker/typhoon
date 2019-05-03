import h from '/utils/h';

import { Input, Checkbox } from '/components/Input';

const Variables = ({ project, questions }) => {
  if (!project) {
    return null;
  }

  if (questions.length === 0) {
    return null;
  }

  return (
    <>
      <h2 style={{ fontWeight: 600, padding: '1em', fontSize: '1.2em' }}>
        Paramètres du ton site{' '}
        <span style={{ fontWeight: 500, fontSize: '0.8em' }}>
          {"(aucun n'est obligatoire, tu peux les laisser vide)"}
        </span>
      </h2>

      {questions.map(question => (
        <Input key={question.name} askIfEmpty {...question} />
      ))}

      <Checkbox title="HTTPS" name="use_https" defaultChecked={project ? project.use_https : false} value="https" />
    </>
  );
};

export default Variables;
