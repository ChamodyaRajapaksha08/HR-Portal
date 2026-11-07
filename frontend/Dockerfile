# 1. Vue (Vite) ප්‍රොජෙක්ට් එක රන් කරන්න Node.js පරිසරය ලබා ගැනීම
FROM node:20-alpine

# 2. Container එක ඇතුළේ වැඩ කරන ප්‍රධාන ෆෝල්ඩර් එක නම් කිරීම
WORKDIR /app

# 3. package.json ෆයිල්ස් ඇතුළට කොපි කිරීම
COPY package*.json ./

# 4. Vue ප්‍රොජෙක්ට් එකට ඕන කරන node_modules ටික ඉන්ස්ටෝල් කිරීම
RUN npm install

# 5. Frontend ෆෝල්ඩර් එකේ තියෙන ඔක්කොම කෝඩ් ෆයිල්ස් Container එකට කොපි කිරීම
COPY . .

# 6. Assignment එකේ ඉල්ලලා තියෙන විදිහට Port 3000 එක Open කිරීම
EXPOSE 3000

# 7. Vite එක හරහා Frontend එක Port 3000 එකෙන් රන් කරන Command එක
CMD ["npm", "run", "dev", "--", "--port", "3000", "--host"]