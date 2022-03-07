import base64
import gzip
import zipfile
import zlib
from io import BytesIO

# msg = "eNrtV0uP40QQ/ivIJxBJ1C/3Izc2O4u4IMQiLgi1OnY7aa3TNm072dFo/jvVdpzN7DgkMwsSEmguI9f7+6qqKw9JblqTLH97SOzHNsB/D8nBtDYUpiy1y5NlIorMFHZNjTIFy4hYc0GJlIIryZFVPJklweSu0u19bUH/4Ao395W3IMgq3xrnbdhVeVf20pPzstpo50elrNV15Xw7eqm3gwdbmrqxuW7dDj5jStSCIolSijDnnPBZUpTV4cygaU1oR3XOUIpTgeWCU4yximJbg/aTmGDVW2vf7dY2gPhbybFAhEgkqUyj/AN8RcnjLGnNpgHPs+RUkqlr0PAmRkw8wKjXXdtWXrcgsTnIzrMRiImFYgilArIxu7p0fqMDwAIaj7MzHgpny1wf/Q4J9p/AISRtA+TvYpQBsFOwKwC8nN1j7Rexu4nkKR6JJFwhxCTB8sTjX7ZTF4L1T4sb8EwVlRycbLq+qmmiJtI6ohvsZsBWj8he4C3GYVd4u5j+Cci9DU2krs8VL+InIG1jW22yrOpOU4AAlY+6aSpdunUw4V5nweaAgDOlbqouZH3/QR/s7WiaPLOJwOquiemart3qtvpgfW83xILiXQO5x256bm33EC9qn0x1Ydtsq/s+A/2IZVllg/ky+eH7n8/8/RTb9sd+rH41pYNlU4XlsrE+f+82vqvf75pVldtfqjPFG6bsaY6TTClYE0z9T1WkqjCu7IL9l5IlFwgpzMg1sk6cjGvgyE2yZPLi0phMz2203Vlg0Wf3uu6a7adkn4l07nLdAN/Ou57MMerndXDojhT2mlL/5b3OyIJRhoiSWKQSky/Y6xFQnmIm1D+718c4DF8h7uUQF2uYvZ0r73Vu9y6zg+EdpVAUpfO3iqg5Y3d4/oYpMifvEL3jK3ZHGL56Vd16MEVCuOCUS0J5yshLDyZov9ckf+nO+uyOQjfRWIcqs1DTwcFie3KqTXGpgALxt1P5ZTfuFCuCMg6Lj8AKTF9/xhYu7Mb34OYrNZrrvbMHyNDk0zOBEeN0cin/PnTF8CQSGHaWCkakYJhFPLbGe1sO09eVZpIkqAZhMfjZ9yc3VmgR/whfQIVffU1URIex9Buwf0UDvn7uGtt8evIZT8VKIjGH3xxozqRazeHQ5PNVuvqOvSHv3sJYHd/VY29kpYu7bHiMhwqH5dWvWvtHsqT88U8l806a"

# r = base64.b64decode(msg)
r = open(r"C:\Users\Administrator\Desktop\1.gzip","rb").read()
dd = zlib.decompress(r)
print(dd.decode("utf8"))
