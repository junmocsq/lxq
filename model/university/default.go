package university

type Classroom struct {
	Building   string `gorm:"primaryKey;size:15"`
	RoomNumber string `gorm:"primaryKey;size:7"`
	Capacity   int
}

func (Classroom) TableName() string {
	return "classroom"
}

type Department struct {
	DeptName string  `gorm:"primaryKey;size:20;comment:院系名称"`
	Building string  `gorm:"size:15;comment:地址"`
	Budget   float32 `gorm:"type:decimal(12,2);check:budget>0;comment:预算"`
}

func (Department) TableName() string {
	return "department"
}

type Course struct {
	CourseId   string     `gorm:"primaryKey;size:7;comment:课程id"`
	Title      string     `gorm:"size:50;comment:课程名称"`
	DeptName   string     `gorm:"index;size:20;comment:院系名称"`
	Credits    float64    `gorm:"decimal;comment:学分"`
	Department Department `gorm:"foreignKey:DeptName;constraint:OnDelete:SET NULL"`
}

func (Course) TableName() string {
	return "course"
}

// Instructor 教师表
type Instructor struct {
	Id         string     `gorm:"size:5;primaryKey;comment:教师id"`
	Name       string     `gorm:"size:20;not null;comment:教师名称"`
	DeptName   string     `gorm:"size:20;comment:院系"`
	Salary     float32    `gorm:"type:decimal(8,2);comment:工资"`
	Department Department `gorm:"foreignKey:DeptName;constraint:OnDelete:SET NULL"`
}

func (Instructor) TableName() string {
	return "instructor"
}

// Section 授课表
type Section struct {
	CourseId   string    `gorm:"size:8;primaryKey;comment:课程id"`
	SecId      string    `gorm:"size:8;primaryKey;comment:预算"`
	Semester   string    `gorm:"size:6;primaryKey;comment:春秋学期"`
	Year       int       `gorm:"primaryKey;comment:学年"`
	Building   string    `gorm:"size:15;comment:建筑"`
	RoomNumber string    `gorm:"size:7;comment:教室号码"`
	TimeSlotId string    `gorm:"size:4;comment:上课时间区间di"`
	Course     Course    `gorm:"foreignKey:CourseId;constraint:OnDelete:cascade"`
	Classroom  Classroom `gorm:"foreignKey:Building,RoomNumber;constraint:OnDelete:SET NULL"`
}

func (Section) TableName() string {
	return "section"
}

type Teaches struct {
	Id         string     `gorm:"size:5;primaryKey;index:mm_idx,priority:5"`
	CourseId   string     `gorm:"size:8;primaryKey;index:mm_idx,priority:2"`
	SecId      string     `gorm:"size:8;primaryKey;index:mm_idx,priority:4"`
	Semester   string     `gorm:"size:6;primaryKey;index:mm_idx,priority:3"`
	Year       int        `gorm:"primaryKey;index:mm_idx,priority:1"`
	Section    Section    `gorm:"foreignKey:CourseId,SecId,Semester,Year;constraint:OnDelete:cascade"`
	Instructor Instructor `gorm:"foreignKey:Id;constraint:OnDelete:cascade"`
}

func (Teaches) TableName() string {
	return "teaches"
}

type Student struct {
	Id         string     `gorm:"primaryKey;size:5"`
	Name       string     `gorm:"size:20;not null"`
	DeptName   string     `gorm:"size:20"`
	TotCred    int        `gorm:"check:tot_cred>=0"`
	Department Department `gorm:"foreignKey:DeptName;constraint:OnDelete:SET NULL"`
}

func (Student) TableName() string {
	return "student"
}

type Takes struct {
	Id       string  `gorm:"primaryKey;size:5"`
	CourseId string  `gorm:"primaryKey;size:8"`
	SecId    string  `gorm:"primaryKey;size:8"`
	Semester string  `gorm:"primaryKey:size:6"`
	Year     int     `gorm:"primaryKey"`
	Grade    string  `gorm:"size:2"`
	Section  Section `gorm:"foreignKey:CourseId,SecId,Semester,Year;constraint:OnDelete:CASCADE"`
	Student  Student `gorm:"foreignKey:Id;constraint:OnDelete:CASCADE"`
}

func (Takes) TableName() string {
	return "takes"
}

type Advisor struct {
	SId        string     `gorm:"primaryKey;size:5"`
	IId        string     `gorm:"size:5"`
	Instructor Instructor `gorm:"foreignKey:IId;references:Id;constraint:OnDelete:SET NULL"`
	Student    Student    `gorm:"foreignKey:SId;references:Id;constraint:OnDelete:CASCADE"`
}

func (Advisor) TableName() string {
	return "advisor"
}

type TimeSlot struct {
	TimeSlotId string `gorm:"primaryKey;size:4"`
	Day        string `gorm:"size:1;primaryKey"`
	StartHr    int    `gorm:"primaryKey;check:end_hr >= 0 and end_hr < 24"`
	StartMin   int    `gorm:"primaryKey;check:end_hr >= 0 and end_hr < 60"`
	EndHr      int    `gorm:"check:end_hr >= 0 and end_hr < 24"`
	EndMin     int    `gorm:"check:end_hr >= 0 and end_hr < 60"`
}

func (TimeSlot) TableName() string {
	return "time_slot"
}

type Prereq struct {
	CourseId string `gorm:"primaryKey;size:8"`
	PrereqId string `gorm:"primaryKey;size:8"`
	Course   Course `gorm:"foreignKey:CourseId;constraint:OnDelete:CASCADE"`
	Course2  Course `gorm:"foreignKey:PrereqId;references:course_id"`
}

func (Prereq) TableName() string {
	return "prereq"
}
