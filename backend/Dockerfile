# 1. Go රන් කරන්න ඕන කරන නිල මෘදුකාංග පරිසරය (Base Image) ලබා ගැනීම
FROM golang:1.25-alpine

# 2. Container එක ඇතුළේ වැඩ කරන ප්‍රධාන ෆෝල්ඩර් එක නම් කිරීම
WORKDIR /app

# 3. Go Packages බාගත කරගන්න go.mod සහ go.sum ෆයිල් ඇතුළට කොපි කිරීම
COPY go.mod go.sum ./
RUN go mod download

# 4. Backend ෆෝල්ඩර් එකේ තියෙන අනෙක් ඔක්කොම කෝඩ් ෆයිල්ස් Container එකට කොපි කිරීම
COPY . .
RUN go mod tidy
# 5. Go කෝඩ් එක රන් කරන්න පුළුවන් Binary ෆයිල් එකක් විදිහට Build කිරීම
RUN go build -o main .

# 6. සර්වර් එක වැඩ කරන Port එක Open කිරීම
EXPOSE 8080

# 7. Container එක ස්ටාර්ට් වෙද්දීම සර්වර් එක ඉබේම රන් කරන Command එක
CMD ["./main"]