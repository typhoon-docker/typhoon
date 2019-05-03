import h from '/utils/h';

import { Redirect } from '@reach/router';

import Projects from '/containers/Projects/';
import { useIsAdmin } from '/utils/connect';

const Admin = () => {
  const isAdmin = useIsAdmin();
  if (!isAdmin) {
    return <Redirect to="/" />;
  }
  return <Projects all />;
};

export default Admin;
