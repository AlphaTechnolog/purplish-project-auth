-- all passwords are 'alpha123.';
INSERT INTO users (id, name, surname, email, local_scopes, company_id, password) VALUES
    ('92ab177b-469d-46c5-b4c2-0e2526ec0710', 'John', 'Doe', 'john.doe@gmail.com', 'r:kardex', 'cde98fd2-6023-4953-9e95-884aed1f09ce', '$2b$10$MSpYfFzZmeVzYEaKKrEVce.UeBIaLx6Xfx1mjOpqgTfOe2JGUPTUK'),
    ('8bfac3c9-d5d5-4559-9c32-fffcb5c6f33b', 'Jane', 'Doe', 'jane.doe@gmail.com', '', 'cde98fd2-6023-4953-9e95-884aed1f09ce', '$2b$10$MSpYfFzZmeVzYEaKKrEVce.UeBIaLx6Xfx1mjOpqgTfOe2JGUPTUK'),
    ('db5ab4f6-1fdb-4b35-8b51-b550105c20cf', 'Gabriel', 'Guerra', 'gfranklings@gmail.com', '*:*', 'b918deaf-92ab-485d-9a69-ee7a2a5f4aef', '$2b$10$MSpYfFzZmeVzYEaKKrEVce.UeBIaLx6Xfx1mjOpqgTfOe2JGUPTUK');

-- users resume
-- -> john.doe@gmail.com:alpha123. -> secondary company, additional scopes: read kardex (company membership: free tier)
-- -> jane.doe@gmail.com:alpha123. -> secondary company, additional scopes: none (company membership: free tier)
-- -> gfranklings@gmail.com:alpha123. -> primary company, additional scopes: everything:admin-like (company membership: premium)